package bundles

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"

	"gopkg.in/yaml.v2"
)

type OrderedJSONElement struct {
	Key, Value interface{}
	index      uint64
}

type OrderedJSON []OrderedJSONElement

func (oj OrderedJSON) Len() int           { return len(oj) }
func (oj OrderedJSON) Less(i, j int) bool { return oj[i].index < oj[j].index }
func (oj OrderedJSON) Swap(i, j int)      { oj[i], oj[j] = oj[j], oj[i] }

var indexCounter uint64

func nextIndex() uint64 {
	indexCounter++
	return indexCounter
}

func (oj OrderedJSON) MarshalJSON() ([]byte, error) {
	buf := &bytes.Buffer{}
	buf.Write([]byte{'{'})
	for i, mi := range oj {
		b, err := json.Marshal(&mi.Value)
		if err != nil {
			return nil, err
		}
		buf.WriteString(fmt.Sprintf("%q:", fmt.Sprintf("%v", mi.Key)))
		buf.Write(b)
		if i < len(oj)-1 {
			buf.Write([]byte{','})
		}
	}
	buf.Write([]byte{'}'})
	return buf.Bytes(), nil
}

func (oj *OrderedJSON) UnmarshalJSON(b []byte) error {
	m := map[string]OrderedJSONElement{}
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}
	for k, v := range m {
		*oj = append(*oj, OrderedJSONElement{Key: k, Value: v.Value, index: v.index})
	}
	sort.Sort(*oj)
	// zero out index for equality checks (its no longer useful and just breaks compares)
	for idx := range *oj {
		(*oj)[idx].index = 0
	}
	return nil
}

func (oje *OrderedJSONElement) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	// if the Unmarshal produced anything but a scalar we may need to convert
	// it to an OrderedJSON object to keep it ordered. It's inefficient to Unmarshal twice but this
	// is the easiest way to inspect the unmarshaled type, and react properly
	up, err := upconvert(b, v)
	if err != nil {
		return nil
	}
	oje.Value = up
	oje.index = nextIndex()
	return nil
}

func (oj *OrderedJSON) UnmarshalYAML(unmarshal func(interface{}) error) error {
	y := yaml.MapSlice{}
	err := unmarshal(&y)
	if err != nil {
		println(err)
	}

	err = oj.ingestYamlMapSlice(&y)
	return err
}

func (oj *OrderedJSON) ingestYamlMapSlice(y *yaml.MapSlice) error {
	for _, ymi := range *y {
		oje, err := convertYamlMapItem(&ymi)
		if err != nil {
			return err
		}
		*oj = append(*oj, oje)
	}
	return nil
}

func convertYamlMapItem(y *yaml.MapItem) (OrderedJSONElement, error) {
	oje := OrderedJSONElement{}
	var err error

	oje.Key = y.Key

	if reflect.TypeOf(y.Value) == reflect.TypeOf(yaml.MapSlice{}) {
		yms := y.Value.(yaml.MapSlice)
		oj := OrderedJSON{}
		oj.ingestYamlMapSlice(&yms)
		oje.Value = oj
		if err != nil {
			return oje, err
		}
	} else if reflect.TypeOf(y.Value) == reflect.TypeOf([]interface{}{}) {
		arr := y.Value.([]interface{})
		for i, v := range arr {
			j := reflect.TypeOf(v)
			_ = j
			if reflect.TypeOf(v) == reflect.TypeOf(yaml.MapSlice{}) {
				yms := v.(yaml.MapSlice)
				oj := OrderedJSON{}
				oj.ingestYamlMapSlice(&yms)
				arr[i] = oj
			}
		}
		oje.Value = arr
	} else {
		oje.Value = y.Value
	}
	oje.index = 0

	return oje, err
}

// utlity function for JSON unmarshaling to unconvert any value that could be an OrderedJSON object
func upconvert(b []byte, any interface{}) (interface{}, error) {
	val := getValue(any)
	switch val.Kind() {
	case reflect.Map:
		var oj OrderedJSON
		err := json.Unmarshal(b, &oj)
		return oj, err

	case reflect.Array, reflect.Slice:
		// There might be a map hidden in a list, so unmarshal into an OrderedJSONElement list
		// (instead of an interface list) so we will inspect and upconvert each element
		ojel := []OrderedJSONElement{}
		if err := json.Unmarshal(b, &ojel); err != nil {
			return []interface{}{}, err
		}
		upconvertedList := make([]interface{}, len(ojel))
		for k, v := range ojel {
			upconvertedList[k] = v.Value
		}
		return upconvertedList, nil
	default:
		return any, nil
	}
}
