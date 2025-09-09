package xtui

import (
	c "github.com/kubex-ecosystem/xtui/components"
	t "github.com/kubex-ecosystem/xtui/types"
)

type Config struct{ t.FormConfig }
type FormFields = t.FormFields
type FormField[T any] struct {
	*t.InputObject[t.FormInputObject[T]]
}
type InputField[T any] struct{ *t.InputObject[T] }

func LogViewer(args ...string) error {
	return c.StartTableScreen(nil, nil)
}
func ShowForm(form Config) (map[string]string, error) {
	return c.ShowForm(form.FormConfig)
}

func NewConfig(title string, fields FormFields) Config {
	return Config{FormConfig: t.NewFormConfig(title, fields.Inputs())}
}
func NewInputField[T any](placeholder string, typ string, value T, required bool, minValue int, maxValue int, err string, validation func(string) error) *FormField[T] {
	input := t.NewInputObject[t.FormInputObject[T]](t.NewFormInputObject[T](value))
	return &FormField[T]{
		input,
	}
}
func NewFormFields[T any](title string, fields []*FormField[t.FormInputObject[T]]) FormFields {
	ffs := make([]t.FormInputObject[any], 0)
	for i, f := range fields {
		ffs[i] = t.NewFormInputObject[any](f.GetValue())
	}
	return FormFields{
		Title:  title,
		Fields: ffs,
	}
}
func NewFormModel(config Config) (map[string]string, error) { return c.ShowForm(config.FormConfig) }
