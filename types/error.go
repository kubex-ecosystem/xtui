package types

// Form and Field Error interface and types

type FormError interface {
	Error() string
	ErrorOrNil() error
	FieldError() map[string]string
	FieldsError() map[string]string
}
type formError struct {
	Rule    string
	Message string
}

func (v *formError) Error() string {
	return v.Message
}
func (v *formError) FieldError() map[string]string {
	return map[string]string{v.Rule: v.Message}
}
func (v *formError) FieldsError() map[string]string {
	return map[string]string{v.Rule: v.Message}
}
func (v *formError) ErrorOrNil() error {
	return v
}

var (
	ErrRequired           = &formError{Rule: "Required", Message: "This field is required"}
	ErrInvalidEmail       = &formError{Rule: "InvalidEmail", Message: "This field must be a valid email address"}
	ErrInvalidURL         = &formError{Rule: "InvalidURL", Message: "This field must be a valid URL"}
	ErrInvalidIP          = &formError{Rule: "InvalidIP", Message: "This field must be a valid IP address"}
	ErrInvalidPort        = &formError{Rule: "InvalidPort", Message: "This field must be a valid Port number"}
	ErrInvalidMin         = &formError{Rule: "InvalidMin", Message: "This field must be a minimum of %d"}
	ErrInvalidMax         = &formError{Rule: "InvalidMax", Message: "This field must be a maximum of %d"}
	ErrInvalidMinLen      = &formError{Rule: "InvalidMinLen", Message: "This field must be a minimum length of %d"}
	ErrInvalidMaxLen      = &formError{Rule: "InvalidMaxLen", Message: "This field must be a maximum length of %d"}
	ErrInvalidRegexp      = &formError{Rule: "InvalidRegexp", Message: "This field must match the regular expression %s"}
	ErrInvalidPattern     = &formError{Rule: "InvalidPattern", Message: "This field must match the pattern %s"}
	ErrInvalidCustom      = &formError{Rule: "InvalidCustom", Message: "This field must match the custom rule"}
	ErrInvalidCustomCheck = &formError{Rule: "InvalidCustomCheck", Message: "This field must match the custom check"}
)
