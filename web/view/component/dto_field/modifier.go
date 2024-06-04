package dto_field

type Modifier func(*FieldConfig)

func Label(l string) Modifier {
	return func(fc *FieldConfig) {
		fc.Label = l
	}
}

func Name(n string) Modifier {
	return func(fc *FieldConfig) {
		fc.Name = n
	}
}

func HelpText(h string) Modifier {
	return func(fc *FieldConfig) {
		fc.HelpText = h
	}
}

func Error(e string) Modifier {
	return func(fc *FieldConfig) {
		fc.HelpText = e
		fc.HelpTextColor = TextColorDanger
		fc.FieldColor = FieldColorDanger
		if fc.FaRightIcon == "" {
			fc.FaRightIcon = "fa-exclamation-triangle"
		}
	}
}

// icon class generally starting with "fa-"
func LeftIcon(li string) Modifier {
	return func(fc *FieldConfig) {
		fc.FaLeftIcon = li
	}
}

// icon class generally starting with "fa-"
func RightIcon(ri string) Modifier {
	return func(fc *FieldConfig) {
		fc.FaRightIcon = ri
	}
}
