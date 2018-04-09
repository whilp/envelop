package main

import "fmt"

type Item struct {
	Details *Details `json:"details"`
}

type Details struct {
	Fields   []*Field   `json:"fields"`
	Sections []*Section `json:"sections"`
}

type Section struct {
	Fields []*SectionField `json:"fields"`
}

type SectionField struct {
	Key        string `json:"k"`
	FieldName  string `json:"n"`
	Title      string `json:"t"`
	FieldValue string `json:"v"`
}

func (f *SectionField) Name() string {
	return f.Title
}

func (f *SectionField) Value() string {
	return f.FieldValue
}

type Field struct {
	Designation string `json:"designation"`
	FieldName   string `json:"name"`
	Type        string `json:"type"`
	FieldValue  string `json:"value"`
}

func (f *Field) Name() string {
	return f.FieldName
}

func (f *Field) Value() string {
	return f.FieldValue
}

type value interface {
	Name() string
	Value() string
}

func (i *Item) Find(key string) (string, error) {
	values := []value{}
	for _, field := range i.Details.Fields {
		values = append(values, field)
	}
	for _, section := range i.Details.Sections {
		for _, field := range section.Fields {
			values = append(values, field)
		}
	}

	for _, field := range values {
		if field.Name() == key {
			return field.Value(), nil
		}
	}
	return "", fmt.Errorf("could not find value for key: %s", key)
}
