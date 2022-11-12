package types

import "fmt"

type FormMsg struct {
	ID    string      `json:"id"`
	Title string      `json:"title"`
	Field []FormField `json:"field"`
}

func (a FormMsg) Convert() (map[string]interface{}, interface{}) {
	builder := MsgBuilder{}
	builder.AppendText(a.Title, TextOption{IsBold: true, IsForm: true})
	for _, field := range a.Field {
		if field.Type == FormFieldText {
			builder.AppendText(field.Placeholder, TextOption{
				IsButton:       true,
				ButtonDataName: fmt.Sprintf("%s|%s", a.ID, field.Key),
				ButtonDataVal:  field.Value.(string),
				ButtonDataAct:  "pub"})
		}
	}
	return builder.Message.Content()
}

type FormField struct {
	Type        FormFieldType      `json:"type"`
	Key         string             `json:"key"`
	Value       interface{}        `json:"value"`
	ValueType   FormFieldValueType `json:"value_type"`
	Required    bool               `json:"required"`
	Placeholder string             `json:"placeholder"`
	Option      interface{}        `json:"option"`
}

type FormFieldType string

const (
	FormFieldText     FormFieldType = "text"
	FormFieldPassword FormFieldType = "password"
	FormFieldNumber   FormFieldType = "number"
	FormFieldColor    FormFieldType = "color"
	FormFieldFile     FormFieldType = "file"
	FormFieldMonth    FormFieldType = "month"
	FormFieldDate     FormFieldType = "date"
	FormFieldTime     FormFieldType = "time"
	FormFieldEmail    FormFieldType = "email"
	FormFieldUrl      FormFieldType = "url"
	FormFieldRadio    FormFieldType = "radio"
	FormFieldCheckbox FormFieldType = "checkbox"
	FormFieldRange    FormFieldType = "range"
	FormFieldSelect   FormFieldType = "select"
	FormFieldTextarea FormFieldType = "textarea"
)

type FormFieldValueType string

const (
	FormFieldValueString       FormFieldValueType = "string"
	FormFieldValueBool         FormFieldValueType = "bool"
	FormFieldValueInt64        FormFieldValueType = "int64"
	FormFieldValueFloat64      FormFieldValueType = "float64"
	FormFieldValueStringSlice  FormFieldValueType = "string_slice"
	FormFieldValueInt64Slice   FormFieldValueType = "int64_slice"
	FormFieldValueFloat64Slice FormFieldValueType = "float64_slice"
)
