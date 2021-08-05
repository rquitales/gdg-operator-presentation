package crd

type CRD struct {
	Metadata struct {
		Name string `json:"name"`
	} `json:"metadata"`
	Spec struct {
		Group string `json:"group"`
		Names struct {
			Kind string `json:"kind"`
		} `json:"names"`
		Versions []struct {
			Name   string `json:"name"`
			Schema struct {
				OpenAPIV3Schema struct {
					Description string `json:"description"`
					Properties  struct {
						Spec struct {
							Description string                       `json:"description"`
							Properties  map[string]map[string]string `json:"properties"`
							Type        string                       `json:"type"`
						} `json:"spec"`
					} `json:"properties"`
					Type string `json:"type"`
				} `json:"openAPIV3Schema"`
			} `json:"schema"`
		} `json:"versions"`
	} `json:"spec"`
}
