package generator

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/manifoldco/promptui"
)

var bundleNameFormat = regexp.MustCompile(`^[a-z0-9-]{5,}`)

var prompts = []func(t *TemplateData) error{
	getSlug,
	getName,
	getAccessLevel,
	getProvisioner,
	getDescription,
}

//TODO: Error Handling
func RunPrompt(t *TemplateData) error {
	var err error
	fmt.Println("in run prompt")

	for _, prompt := range prompts {
		err = prompt(t)
		if err != nil {
			return err
		}
	}

	return nil
}

func getSlug(t *TemplateData) error {
	validate := func(input string) error {
		if !bundleNameFormat.MatchString(input) {
			return errors.New("name must be greater than 4 characters and can only include lowercase letters and dashes")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Slug",
		Validate: validate,
	}

	result, err := prompt.Run()
	if err != nil {
		return err
	}

	t.Slug = result
	return nil
}

func getName(t *TemplateData) error {
	prompt := promptui.Prompt{
		Label: "Name",
	}

	result, err := prompt.Run()

	if err != nil {
		return err
	}

	t.Name = result
	return nil
}

func getAccessLevel(t *TemplateData) error {
	prompt := promptui.Select{
		Label: "Access Level",
		Items: []string{"Public", "Private"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		return err
	}

	t.Access = result
	return nil
}

func getProvisioner(t *TemplateData) error {
	prompt := promptui.Select{
		Label: "Provisioner",
		Items: []string{"terraform"},
	}
	_, result, err := prompt.Run()

	if err != nil {
		return err
	}

	t.Provisioner = result
	return nil
}

func getDescription(t *TemplateData) error {
	prompt := promptui.Prompt{
		Label: "Description",
	}

	result, err := prompt.Run()

	if err != nil {
		return err
	}

	t.Description = result
	return nil
}
