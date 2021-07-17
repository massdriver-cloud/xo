package transform

import (
	"io/ioutil"

	"github.com/robertkrimen/otto"
)

func Transform(input string, transformerFile string) (string, error) {
	transformer, err := ioutil.ReadFile(transformerFile)
	if err != nil {
		return "", err
	}

	// load the transform function into the VM
	vm := otto.New()
	_, err = vm.Run(string(transformer))
	if err != nil {
		return "", err
	}

	// call the transform function
	out, err := vm.Call("transform", nil, input)
	if err != nil {
		return "", err
	}

	return out.ToString()
}
