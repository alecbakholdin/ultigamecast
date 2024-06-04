package dto_field

type InputType string

const (
	InputTypeText     InputType = "text"
	InputTypeEmail    InputType = "email"
	InputTypePassword InputType = "password"
)

func TextType(c *FieldConfig) {
	c.Type = InputTypeText
}

func EmailType(c *FieldConfig) {
	c.Type = InputTypeEmail
}

func PasswordType(c *FieldConfig) {
	c.Type = InputTypePassword
}