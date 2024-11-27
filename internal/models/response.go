package models

type SuccessResponse struct {
	Result   interface{} `json:"result"`
	Metadata Properties  `json:"metadata"`
}

type ErrorResponse struct {
	Result []ErrorDetail `json:"errors"`
}

type ErrorDetail struct {
	Code     string     `json:"code"`
	Message  string     `json:"message"`
	Metadata Properties `json:"metadata"`
}

type Properties struct {
	Properties1 string `json:"additionalProp1,omitempty"`
	Properties2 string `json:"additionalProp2,omitempty"`
	Properties3 string `json:"additionalProp3,omitempty"`
}
