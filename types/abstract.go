// Package types provides type definitions and interfaces for form handling and validation.
package types

import (
	"github.com/charmbracelet/lipgloss"
)

// Validation interfaces and types

type FormInputValidationObject[T any] func(value T) error
type FormInputValidationString func(value string) error
type FormInputValidationCustom func(value string) error
type FormInputValidationRule[T any] struct {
	Min, Max  int
	ObjFunc   FormInputValidationObject[T]
	Challenge FormInputValidationCustom
	StrFunc   FormInputValidationString
}

// ---------------------------------------------------
// Form objects interfaces and types for handle with groups of fields, forms, sections and parts.

// FormPart is an interface that agroups fields and provides methods to manipulate them. It is used to
// group fields, manage screen space and fields distribution on the screen.
type FormPart struct {
	*lipgloss.Style
	FormGroup
	Width, Height, MaxWidth, MaxHeight int
	Title                              string
}

func (s FormPart) GetWidth() int { return s.Width }
func (s FormPart) GetUpperBound() int {
	mh := s.GetMaxHeight()
	if mh == 0 {
		return 0
	}
	return s.GetLowerBound() + mh
}
func (s FormPart) GetLowerBound() int {
	if s.GetMaxHeight() == 0 {
		return 0
	}
	return s.GetMaxHeight() - s.GetHeight()
}
func (s FormPart) GetLeftBound() int {
	if s.GetMaxWidth() == 0 {
		return 0
	}
	return s.GetMaxWidth() - s.GetWidth()
}
func (s FormPart) GetRightBound() int {
	mw := s.GetMaxWidth()
	if mw == 0 {
		return 0
	}
	return s.GetLeftBound() + mw
}

// FormGroup is an interface that manages goups of fields. It was created to provide a easy way to
// manage fields in a form, wrapping them in groups and providing methods to manipulate them.
type FormGroup interface {
	GetFields() FormFields

	GetFieldByID(id string) FieldDefinition
	GetFieldByIndex(index int) FieldDefinition
	GetFieldIndex(id string) int
	GetFieldID(index int) string

	FieldsCount() int

	Validate() error

	SetField(index int, field FieldDefinition)
	SetFieldByID(id string, field FieldDefinition)
	SetFields(fields FormFields)
}
