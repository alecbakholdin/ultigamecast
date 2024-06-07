package mod

import "ultigamecast/web/view/component/dto/field"

type HelpModifier struct {
	HelpTextColor field.TextColor
	HelpText string

	TooltipHelpText string
}

func (h *HelpModifier) Apply(fc *field.FieldConfig) {
	if h.HelpText != "" {
		fc.HelpText = h.HelpText
	}

	if h.HelpTextColor != "" {
		fc.HelpTextColor = h.HelpTextColor
	}

	if h.TooltipHelpText != "" {
		fc.TooltipHelpText = h.TooltipHelpText
	}
}

func HelpText(h string) *HelpModifier{
	return &HelpModifier{
		HelpText: h,
	}
}

func HelpTextColor(c field.TextColor) *HelpModifier {
	return &HelpModifier{
		HelpTextColor: c,
	}
}

func TooltipHelpText(h string) *HelpModifier {
	return &HelpModifier{
		TooltipHelpText: h,
	}
}