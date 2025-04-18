package model

type Details struct {
	ParentModel       string   `json:"parent_model"`
	Format            string   `json:"format"`
	Family            string   `json:"family"`
	Families          []string `json:"families"`
	ParameterSize     string   `json:"parameter_size"`
	QuantizationLevel string   `json:"quantization_level"`
}

type Model struct {
	Name       string  `json:"name"`
	Model      string  `json:"model,omitempty"`
	ModifiedAt string  `json:"modified_at,omitempty"`
	Size       int64   `json:"size,omitempty"`
	Digest     string  `json:"digest,omitempty"`
	Details    Details `json:"details,omitempty"`
}
