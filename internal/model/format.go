package model 

type Format struct {
	Type string `json:type`
	Qualities []string `json:qualities`
}

type FormatResponse struct {
	Formats []Format `json:formats`
}