package generator

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

var prompts = []func(t *TemplateData){
	getName,
	getProvisioner,
	getDescription,
}

//TODO: Validation and Error Handling
func RunPrompt(t *TemplateData) *TemplateData {
	fmt.Println("in run prompt")

	for _, prompt := range prompts {
		prompt(t)
	}

	return t
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
