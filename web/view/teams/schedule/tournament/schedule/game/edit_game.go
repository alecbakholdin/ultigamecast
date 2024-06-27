package view_game

import "ultigamecast/web/view/component/dto"

type EditGameDTO struct {
	dto.DTO
	Field string `validate:"required,ascii"`
	Value string `validate:"ascii"`
}