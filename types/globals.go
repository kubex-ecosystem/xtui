package types

// Field Basic Generic Definition Interface

type FieldDefinition interface {
	Description() string
	String() string
}

// Shortcut Types

type Shortcut string

const (
	Q Shortcut = "q"
	F Shortcut = "f"

	UP    Shortcut = "up"
	DOWN  Shortcut = "down"
	LEFT  Shortcut = "left"
	RIGHT Shortcut = "right"

	DEL   Shortcut = "del"
	TAB   Shortcut = "tab"
	ESC   Shortcut = "esc"
	ENTER Shortcut = "enter"

	CTRLA Shortcut = "ctrl+a"
	CTRLE Shortcut = "ctrl+e"
	CTRLR Shortcut = "ctrl+r"
	CTRLC Shortcut = "ctrl+c"
	CTRLH Shortcut = "ctrl+h"

	SHIFTTAB Shortcut = "shift+tab"
)

func (s Shortcut) Description() string { return "Shortcut " + string(s) }
func (s Shortcut) String() string      { return string(s) }

// Field Position and Size Types

type FieldSize string

const (
	SizeDefault FieldSize = "default"
	SizeSmall   FieldSize = "small"
	SizeLarge   FieldSize = "large"
)

func (f FieldSize) Description() string { return "Field Size " + string(f) }
func (f FieldSize) String() string      { return string(f) }

type FieldPosition string

const (
	PositionDefault FieldPosition = "default"
	PositionTop     FieldPosition = "top"
	PositionBottom  FieldPosition = "bottom"
)

func (f FieldPosition) Description() string { return "Field Position " + string(f) }
func (f FieldPosition) String() string      { return string(f) }

// Field Alignment Types

type FieldAlignment string

const (
	AlignmentDefault FieldAlignment = "default"
	AlignmentLeft    FieldAlignment = "left"
	AlignmentCenter  FieldAlignment = "center"
	AlignmentRight   FieldAlignment = "right"
)

func (f FieldAlignment) Description() string { return "Field Alignment " + string(f) }
func (f FieldAlignment) String() string      { return string(f) }

// Field Input Primitive/Generic Types

type FieldType string

const (
	FieldBool     FieldType = "bool"
	FieldInt      FieldType = "int"
	FieldText     FieldType = "text"
	FieldPass     FieldType = "password"
	FieldDate     FieldType = "date"
	FieldTime     FieldType = "time"
	FieldList     FieldType = "list"
	FieldFile     FieldType = "file"
	FieldTable    FieldType = "table"
	FieldFunction FieldType = "function"
)

func (f FieldType) Description() string { return "Field Type " + string(f) }
func (f FieldType) String() string      { return string(f) }

// Field Rules and Validation Types

type FieldRule interface {
	Validate(value string) error
}

type ValidationRule string

const (
	Required ValidationRule = "required"
	Email    ValidationRule = "email"
	URL      ValidationRule = "url"
	IP       ValidationRule = "ip"
	Port     ValidationRule = "port"
	Min      ValidationRule = "min"
	Max      ValidationRule = "max"
	MinLen   ValidationRule = "min_len"
	MaxLen   ValidationRule = "max_len"
	Regexp   ValidationRule = "regexp"
	Pattern  ValidationRule = "pattern"
)

func (v ValidationRule) Description() string { return "Validation Rule " + string(v) }
func (v ValidationRule) String() string      { return string(v) }
func (v ValidationRule) Validate(value string, customCheck func(interface{}) error) error {
	switch v {
	case Required:
		if value == "" {
			return ErrRequired
		}
		// TODO: Add more native basic validation rules
		//default:
		//	if customCheck != nil {
		//		return customCheck(v)
		//	}
	}
	if customCheck != nil {
		return customCheck(v)
	}
	return nil
}
