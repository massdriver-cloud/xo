package generator

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/manifoldco/promptui"
)

var bundleNameFormat = regexp.MustCompile(`^[a-z0-9-]{5,}`)

var prompts = []func(t *TemplateData){
	getSlug,
	getName,
	getAccessLevel,
	getProvisioner,
	getDescription,
}

//TODO: Error Handling
func RunPrompt(t *TemplateData) *TemplateData {
	fmt.Println("in run prompt")

	for _, prompt := range prompts {
		prompt(t)
	}

	return t
}

func getSlug(t *TemplateData) {
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
		fmt.Printf("Error: %s", err)
	}

	t.Slug = result
}

func getName(t *TemplateData) {
	prompt := promptui.Prompt{
		Label: "Name",
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Error: %s", err)
	}

	t.Name = result
}

func getAccessLevel(t *TemplateData) {
	prompt := promptui.Select{
		Label: "Access Level",
		Items: []string{"Public", "Private"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Error: %s", err)
	}

	t.Access = result
}

func getProvisioner(t *TemplateData) {
	prompt := promptui.Select{
		Label: "Provisioner",
		Items: []string{"terraform"},
	}
	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Error: %s", err)
	}

	t.Provisioner = result
}

func getDescription(t *TemplateData) {
	prompt := promptui.Prompt{
		Label: "Description",
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Error: %s", err)
	}

	t.Description = result
}
