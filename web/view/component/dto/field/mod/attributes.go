package mod

import "ultigamecast/web/view/component/dto/field"

type AttributeModifier struct {
	Key   string
	Value string
}

func (a *AttributeModifier) Apply(fc *field.FieldConfig) {
	fc.InputAttributes[a.Key] = a.Value
}

func InputAttribute(k, v string) *AttributeModifier {
	return &AttributeModifier{
		Key:   k,
		Value: v,
	}
}

func Name(n string) *AttributeModifier {
	return InputAttribute("name", n)
}

func Placeholder(p string) *AttributeModifier {
	return InputAttribute("placeholder", p)
}

func Autocomplete(a string) *AttributeModifier {
	return InputAttribute("autocomplete", a)
}

func InputType(t field.InputType) *AttributeModifier {
	return InputAttribute("type", string(t))
}

func InputTypeText() *AttributeModifier {
	return InputType(field.InputTypeText)
}

func InputTypeEmail() *AttributeModifier {
	return InputType(field.InputTypeEmail)
}

func InputTypePassword() *AttributeModifier {
	return InputType(field.InputTypePassword)
}
