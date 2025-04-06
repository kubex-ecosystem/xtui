package types

type Config struct {
	Title  string
	Fields FormFields
}

func (c Config) GetTitle() string      { return c.Title }
func (c Config) GetFields() FormFields { return c.Fields }

type FormConfig struct {
	Title string
	FormFields
}

func NewFormConfig(title string, fields []FormInputObject[any]) FormConfig {
	return FormConfig{
		Title:      title,
		FormFields: FormFields{Title: title, Fields: fields},
	}
}

func (f FormConfig) GetTitle() string      { return f.Title }
func (f FormConfig) SetTitle(title string) { f.Title = title }

func (f FormConfig) GetFields() []FormInputObject[any] { return f.Fields }
func (f FormConfig) GetField(name string) FormInputObject[any] {
	for _, field := range f.Fields {
		if field.GetName() == name {
			return field
		}
	}
	return nil
}
func (f FormConfig) GetFieldValue(name string) any {
	for _, field := range f.Fields {
		if field.GetName() == name {
			return field.GetValue()
		}
	}
	return nil
}

func (f FormConfig) SetFieldValue(name string, value any) error {
	for _, field := range f.Fields {
		if field.GetName() == name {
			return field.SetValue(value)
		}
	}
	return nil
}
func (f FormConfig) SetField(name string, field FormInputObject[any]) error {
	for i, cf := range f.GetFields() {
		if cf.GetName() == name {
			f.Fields[i] = field
			return nil
		}
	}
	return nil
}

func (f FormConfig) AddField(field FormInputObject[any]) { f.Fields = append(f.Fields, field) }
func (f FormConfig) RemoveField(name string) error {
	for i, field := range f.Fields {
		if field.GetName() == name {
			f.Fields = append(f.Fields[:i], f.Fields[i+1:]...)
			return nil
		}
	}
	return nil
}
