package models

type Data struct {
	Stats     Stats       `json:"stats"`
	Escaped   []Entry     `json:"escaped"`
	Timeouted interface{} `json:"timeouted"`
	Killed    []Entry     `json:"killed"`
	Errored   interface{} `json:"errored"`
}
