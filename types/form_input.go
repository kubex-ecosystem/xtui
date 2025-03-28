package types

import (
	"fmt"
	"reflect"
)

// FormInputObject is the most basic form input object. It will be used globally in all form input objects.
type FormInputObject[T any] interface {
	GetType() reflect.Type
	GetValue() T
	SetValue(val T) error
}

// InputObject is a struct that implements the FormInputObject interface. It is used to store the value
// of the form input object in a serializable. It is used to store the value of the form input object in a
// serializable and to store metadata for easy and integrated serialization and type conversion.
type InputObject[T any] struct {
	Val T `json:"value" yaml:"value" gorm:"column:value"`
	err error
}

func (s *InputObject[T]) GetType() reflect.Type { return reflect.TypeOf(s.Val) }
func (s *InputObject[T]) GetValue() T           { return s.Val }
func (s *InputObject[T]) SetValue(val T) error {
	if s != nil {
		s.Val = val
		return nil
	} else {
		return fmt.Errorf("input is nil")
	}
}

// -----------------------------------------------------------------------------
// Here we have that will be used to create another package later, for now it is just a draft.

type FormInput[T any] interface {
	FieldDefinition
	FormInputObject[T]

	// Common Getters

	Placeholder() string
	MinValue() int
	MaxValue() int
	Validation() func(string, func(interface{}) error) error

	// Boolean methods

	IsRequired() bool
	Error() string

	// Common Setters

	SetPlaceholder(string)

	// Validation methods

	SetRequired(bool)
	SetMinValue(int)
	SetMaxValue(int)
	SetValidation(func(string, func(interface{}) error) error)
	SetValidationRules([]ValidationRule)
	ValidationRules() []ValidationRule
	Validate() error

	// Factory methods

	String() string
	FromString(string) error
	ToMap() map[string]interface{}
	FromMap(map[string]interface{}) error
}

type Input[T any] struct {
	FieldDefinition
	FormInputObject[T]
	Ph                 string           `json:"placeholder" yaml:"placeholder" gorm:"column:placeholder"`
	Tp                 reflect.Type     `json:"type" yaml:"type" gorm:"column:type"`
	Val                *T               `json:"value" yaml:"value" gorm:"column:value"`
	Req                bool             `json:"required" yaml:"required" gorm:"column:required"`
	Min                int              `json:"min" yaml:"min" gorm:"column:min"`
	Max                int              `json:"max" yaml:"max" gorm:"column:max"`
	Err                string           `json:"error" yaml:"error" gorm:"column:error"`
	ValidationRulesVal []ValidationRule `json:"validation_rules" yaml:"validation_rules" gorm:"column:validation_rules"`
}

func (s *Input[T]) GetType() reflect.Type { return reflect.TypeOf(s.Val) }
func (s *Input[T]) GetValue() T {
	v := reflect.ValueOf(s.Val)
	return v.Interface().(T)
}
func (s *Input[T]) SetValue(val T) error {
	if s != nil {
		s.Val = &val
		return nil
	} else {
		return fmt.Errorf("input is nil")
	}
}

func (s *Input[T]) MinValue() int                                           { return s.Min }
func (s *Input[T]) MaxValue() int                                           { return s.Max }
func (s *Input[T]) Validation() func(string, func(interface{}) error) error { return nil }
func (s *Input[T]) Error() string                                           { return s.Err }
func (s *Input[T]) SetMaxValue(i int)                                       { s.Max = i }
func (s *Input[T]) Placeholder() string                                     { return s.Ph }

// // Boolean methods

func (s *Input[T]) IsRequired() bool { return s.Req }
func (s *Input[T]) GetError() string { return s.Err }

// // Common Setters

func (s *Input[T]) SetPlaceholder(ph string) { s.Ph = ph }

// // Validation methods

func (s *Input[T]) SetRequired(req bool) { s.Req = req }
func (s *Input[T]) SetMinValue(min int)  { s.Min = min }
func (s *Input[T]) SetValidation(validation func(string, func(interface{}) error) error) {
	var valRules = make([]ValidationRule, 0)
	if s != nil {
		if s.ValidationRules() != nil {
			if len(s.ValidationRules()) >= 1 {
				valRule := s.ValidationRules()[0]
				if err := valRule.Validate(s.String(), nil); err != nil {
					s.Err = err.Error()
				}
			}
		}
	}
	s.ValidationRulesVal = valRules
}
func (s *Input[T]) SetValidationRules(rules []ValidationRule) { s.ValidationRulesVal = rules }
func (s *Input[T]) ValidationRules() []ValidationRule         { return s.ValidationRulesVal }
func (s *Input[T]) Validate() error {
	if s.ValidationRulesVal == nil {
		return nil
	}
	for _, rule := range s.ValidationRulesVal {
		if s.Val != nil {
			if err := rule.Validate(s.String(), nil); err != nil {
				return err
			}
		}
	}
	return nil
}
func (s *Input[T]) String() string {
	if s != nil {
		if s.Val != nil {
			if s.GetType().Kind() == reflect.String {
				v := s.GetValue()
				vv := reflect.ValueOf(v)
				return vv.String()
			}
		}
	}
	return ""
}

// // Factory methods

func (s *Input[T]) FromString(str string) error {
	if s != nil {
		if s.Val != nil {
			v := s.GetValue()
			vv := reflect.ValueOf(v)
			vv.SetString(str)
		}
	}
	return nil
}
func (s *Input[T]) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"placeholder": s.Ph,
		"type":        s.Tp,
		"value":       s.Val,
		"required":    s.Req,
		"min":         s.Min,
		"max":         s.Max,
		"error":       s.Err,
	}
}
func (s *Input[T]) FromMap(m map[string]interface{}) error { return nil }

func NewInput[T FormInputObject[any]](t T) *Input[T]        { return &Input[T]{Val: &t} }
func NewFormInput[T FormInputObject[any]](t T) FormInput[T] { return NewInput[T](t) }
func NewFormInputFromMap[T FormInputObject[any]](m map[string]interface{}) FormInput[T] {
	return NewInput[T](m["value"].(T))
}
func NewFormInputFromString[T FormInputObject[any]](str string) FormInput[T] {
	v := reflect.ValueOf(str)
	return NewInput(v.Interface().(T))
}
func NewFormInputFromBytes[T FormInputObject[any]](b []byte) FormInput[T] {
	v := reflect.ValueOf(b)
	return NewInput[T](v.Interface().(T))
}

func NewInputObject[T any](t T) *InputObject[T]        { return &InputObject[T]{Val: t} }
func NewFormInputObject[T any](t T) FormInputObject[T] { return NewInputObject[T](t) }
func NewFormInputObjectFromMap[T any](m map[string]interface{}) FormInputObject[T] {
	v := m["value"].(T)
	return NewInputObject[T](v)
}
func NewFormInputObjectFromString[T any](str string) FormInputObject[T] {
	v := reflect.ValueOf(str)
	return NewInputObject[T](v.Interface().(T))
}
func NewFormInputObjectFromBytes[T any](b []byte) FormInputObject[T] {
	v := reflect.ValueOf(b)
	return NewInputObject[T](v.Interface().(T))
}
