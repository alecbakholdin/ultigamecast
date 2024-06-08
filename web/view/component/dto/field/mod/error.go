package mod

import "ultigamecast/web/view/component/dto/field"

type IncludeFormErrorModifier struct {

}

func (*IncludeFormErrorModifier) Apply(fc *field.FieldConfig) {
	if fe := fc.Dto.FormError(); fe != "" {
		if fc.HelpText == "" {
			fc.HelpText = fe
		} else {
			fc.HelpText += ", " + fe
		}
		fc.HelpTextColor = field.TextColorDanger
	}
}

func IncludeFormError() *IncludeFormErrorModifier {
	return &IncludeFormErrorModifier{}
}