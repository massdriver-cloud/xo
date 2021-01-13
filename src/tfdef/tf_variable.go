package tfdef

type TFVariableFile struct {
	Variable map[string]TFVariable `json:"variable"`
}

// TFVariable is a representation of a terraform HCL "variable"
type TFVariable struct {
	Type string `json:"type"`
}

// func (tfv TFVariable) New(p Property) Property {}
// func (tfvf TFVariableFile) MarshalJSON() ([]byte, error) {}
