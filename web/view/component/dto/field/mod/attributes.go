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

func Value(v string) *AttributeModifier {
	return InputAttribute("value", v)
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

func Autofocus() *AttributeModifier {
	return InputAttribute("autofocus", "true")
}

func InputType(t field.InputType) *AttributeModifier {
	return InputAttribute("type", string(t))
}

func Id(id string) *AttributeModifier {
	return InputAttribute("id", id)
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

func InputTypeHidden() *AttributeModifier {
	return InputType(field.InputTypeHidden)
}

func InputTypeDatetime() *AttributeModifier {
	return InputType(field.InputTypeDatetime)
}

func InputTypeNumber() *AttributeModifier {
	return InputType(field.InputTypeNumber)
}