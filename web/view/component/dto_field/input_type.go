package dto_field

type InputType string

const (
	InputTypeText     InputType = "text"
	InputTypeEmail    InputType = "email"
	InputTypePassword InputType = "password"
)

func TypeText() Modifier {
	return func(c *FieldConfig) {
		c.Type = InputTypeText
	}
}

func TypeEmail() Modifier {
	return func(c *FieldConfig) {
		c.Type = InputTypeEmail
	}
}

func TypePassword() Modifier {
	return func(c *FieldConfig) {
		c.Type = InputTypePassword
	}
}
