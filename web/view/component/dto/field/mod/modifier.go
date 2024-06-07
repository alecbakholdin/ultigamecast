package mod

import "ultigamecast/web/view/component/dto/field"

type Modifier interface {
	Apply(*field.FieldConfig)
}