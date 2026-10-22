package resource

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"xo/src/bundle"

	"github.com/itchyny/gojq"
)

// Template is a compiled jq template that renders the payload for a single
// resource field.
type Template struct {
	Field string
	// Origin describes where the template came from (inline step block or a
	// file path) and is used in error messages.
	Origin string
	Source string
	code   *gojq.Code
}

// TemplateError is returned when a template fails to parse, compile or
// evaluate. Its message includes the field, where the template came from and,
// when the location is known, the offending line with a caret under it.
type TemplateError struct {
	Field  string
	Origin string
	Stage  string
	Err    error
	// Detail holds an optional source excerpt and hints.
	Detail string
}

func (e *TemplateError) Error() string {
	msg := fmt.Sprintf("resource %q: jq template from %s failed to %s: %v", e.Field, e.Origin, e.Stage, e.Err)
	if e.Detail != "" {
		msg += "\n" + e.Detail
	}
	return msg
}

func (e *TemplateError) Unwrap() error { return e.Err }

// templateFilePrefixes are the template file name prefixes, in priority order.
var templateFilePrefixes = []string{"resource_", "artifact_"}

// LoadTemplates discovers and compiles the jq templates for a step. Inline
// templates take priority over resource_<field>.jq files, which take priority
// over legacy artifact_<field>.jq files in templateDir. Every template is
// compiled before returning so that a broken template is reported before
// anything is published.
func LoadTemplates(inline map[string]bundle.StepResource, templateDir string) ([]*Template, []string, error) {
	byField, warnings, err := discoverTemplates(inline, templateDir)
	if err != nil {
		return nil, warnings, err
	}

	var errs []error
	templates := make([]*Template, 0, len(byField))
	for _, field := range sortedFields(byField) {
		tmpl := byField[field]
		if err := tmpl.compile(); err != nil {
			errs = append(errs, err)
			continue
		}
		templates = append(templates, tmpl)
	}
	if len(errs) > 0 {
		return nil, warnings, errors.Join(errs...)
	}
	return templates, warnings, nil
}

// TemplateFields returns the fields a step has a template for, inline or as a
// file, without compiling them. It is used on decommission to find the
// resources a step manages, where a broken template must not block cleanup.
func TemplateFields(inline map[string]bundle.StepResource, templateDir string) ([]string, error) {
	byField, _, err := discoverTemplates(inline, templateDir)
	if err != nil {
		return nil, err
	}
	return sortedFields(byField), nil
}

func discoverTemplates(inline map[string]bundle.StepResource, templateDir string) (map[string]*Template, []string, error) {
	var warnings []string
	byField := map[string]*Template{}

	for field, res := range inline {
		if res.Template == "" {
			continue
		}
		byField[field] = &Template{Field: field, Origin: fmt.Sprintf("massdriver.yaml steps[].resources.%s.template", field), Source: res.Template}
	}

	for _, prefix := range templateFilePrefixes {
		matches, err := filepath.Glob(filepath.Join(templateDir, prefix+"*.jq"))
		if err != nil {
			return nil, nil, err
		}
		sort.Strings(matches)
		for _, path := range matches {
			field := strings.TrimSuffix(strings.TrimPrefix(filepath.Base(path), prefix), ".jq")
			if field == "" {
				continue
			}
			if existing, ok := byField[field]; ok {
				warnings = append(warnings, fmt.Sprintf("resource %q: ignoring %s, using template from %s", field, path, existing.Origin))
				continue
			}
			data, err := os.ReadFile(path)
			if err != nil {
				return nil, nil, fmt.Errorf("resource %q: unable to read template %s: %w", field, path, err)
			}
			byField[field] = &Template{Field: field, Origin: path, Source: string(data)}
		}
	}

	return byField, warnings, nil
}

func sortedFields(byField map[string]*Template) []string {
	fields := make([]string, 0, len(byField))
	for field := range byField {
		fields = append(fields, field)
	}
	sort.Strings(fields)
	return fields
}

func (t *Template) compile() error {
	if strings.TrimSpace(t.Source) == "" {
		return &TemplateError{Field: t.Field, Origin: t.Origin, Stage: "parse", Err: errors.New("template is empty")}
	}

	query, err := gojq.Parse(t.Source)
	if err != nil {
		tErr := &TemplateError{Field: t.Field, Origin: t.Origin, Stage: "parse", Err: err}
		var pErr *gojq.ParseError
		if errors.As(err, &pErr) {
			offset := pErr.Offset - len(pErr.Token)
			if offset < 0 {
				offset = 0
			}
			line, col, excerpt := locate(t.Source, offset)
			tErr.Err = fmt.Errorf("%w at line %d, column %d", err, line, col)
			tErr.Detail = excerpt
		}
		return tErr
	}

	code, err := gojq.Compile(query)
	if err != nil {
		tErr := &TemplateError{Field: t.Field, Origin: t.Origin, Stage: "compile", Err: err}
		if strings.Contains(err.Error(), "variable not defined") {
			tErr.Detail = "  hint: inputs are available on the input object as .params, .connections, .envs, .secrets and .outputs, not as $variables"
		}
		return tErr
	}
	t.code = code
	return nil
}

// Render evaluates the template against inputs. The template must produce
// exactly one JSON object.
func (t *Template) Render(inputs map[string]any) (map[string]any, error) {
	iter := t.code.Run(inputs)
	var results []any
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, isErr := v.(error); isErr {
			var halt *gojq.HaltError
			if errors.As(err, &halt) && halt.Value() == nil {
				break
			}
			return nil, &TemplateError{Field: t.Field, Origin: t.Origin, Stage: "evaluate", Err: err, Detail: inputHint(inputs)}
		}
		results = append(results, v)
	}

	switch len(results) {
	case 0:
		return nil, &TemplateError{Field: t.Field, Origin: t.Origin, Stage: "evaluate", Err: errors.New("template produced no output (expected exactly one object)")}
	case 1:
	default:
		return nil, &TemplateError{Field: t.Field, Origin: t.Origin, Stage: "evaluate",
			Err:    fmt.Errorf("template produced %d outputs (expected exactly one object)", len(results)),
			Detail: "  hint: an iterator such as .[] or a comma at the top level emits multiple values; wrap it in [...] to collect them into an array"}
	}

	obj, ok := results[0].(map[string]any)
	if !ok {
		return nil, &TemplateError{Field: t.Field, Origin: t.Origin, Stage: "evaluate",
			Err: fmt.Errorf("template produced %s (expected a JSON object)", typeName(results[0]))}
	}
	return obj, nil
}

// NullPaths returns the paths in v whose value is null. A null usually means
// the template referenced an input that does not exist (e.g. a typo in an
// output name), so callers surface these as warnings.
func NullPaths(v any) []string {
	var paths []string
	var walk func(prefix string, v any)
	walk = func(prefix string, v any) {
		switch val := v.(type) {
		case nil:
			paths = append(paths, prefix)
		case map[string]any:
			keys := make([]string, 0, len(val))
			for k := range val {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				walk(prefix+"."+k, val[k])
			}
		case []any:
			for i, item := range val {
				walk(fmt.Sprintf("%s[%d]", prefix, i), item)
			}
		}
	}
	walk("", v)
	return paths
}

// locate converts a byte offset into a 1-based line and column and returns
// an excerpt of that line with a caret under the column.
func locate(src string, offset int) (int, int, string) {
	if offset > len(src) {
		offset = len(src)
	}
	line := strings.Count(src[:offset], "\n") + 1
	lineStart := strings.LastIndex(src[:offset], "\n") + 1
	lineEnd := strings.IndexByte(src[lineStart:], '\n')
	if lineEnd < 0 {
		lineEnd = len(src)
	} else {
		lineEnd += lineStart
	}
	col := offset - lineStart + 1
	text := src[lineStart:lineEnd]
	lineNum := fmt.Sprintf("%d", line)
	// Preserve tabs so the caret lines up with the source line.
	pad := strings.Map(func(r rune) rune {
		if r == '\t' {
			return '\t'
		}
		return ' '
	}, text[:col-1])
	excerpt := fmt.Sprintf("  %s | %s\n  %s | %s^", lineNum, text, strings.Repeat(" ", len(lineNum)), pad)
	return line, col, excerpt
}

func inputHint(inputs map[string]any) string {
	var parts []string
	for _, key := range []string{"params", "connections", "envs", "secrets", "outputs"} {
		switch val := inputs[key].(type) {
		case map[string]any:
			keys := make([]string, 0, len(val))
			for k := range val {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			if len(keys) == 0 {
				parts = append(parts, fmt.Sprintf("  .%s: (empty)", key))
			} else {
				parts = append(parts, fmt.Sprintf("  .%s keys: %s", key, strings.Join(keys, ", ")))
			}
		case []any:
			parts = append(parts, fmt.Sprintf("  .%s: %s", key, describeArray(val)))
		}
	}
	if len(parts) == 0 {
		return ""
	}
	return "  available inputs:\n  " + strings.Join(parts, "\n  ")
}

// describeArray summarizes an array input. When its items carry a "kind"
// (e.g. Kubernetes manifests from Helm), the distinct kinds are listed since
// templates typically select on them.
func describeArray(items []any) string {
	if len(items) == 0 {
		return "(empty array)"
	}
	kindSet := map[string]bool{}
	for _, item := range items {
		if obj, ok := item.(map[string]any); ok {
			if kind, ok := obj["kind"].(string); ok && kind != "" {
				kindSet[kind] = true
			}
		}
	}
	desc := fmt.Sprintf("array of %d items", len(items))
	if len(kindSet) > 0 {
		desc += ", kinds: " + strings.Join(sortedKeys(kindSet), ", ")
	}
	return desc
}

func typeName(v any) string {
	switch v.(type) {
	case nil:
		return "null"
	case bool:
		return "a boolean"
	case string:
		return "a string"
	case []any:
		return "an array"
	case map[string]any:
		return "an object"
	default:
		return "a number"
	}
}
