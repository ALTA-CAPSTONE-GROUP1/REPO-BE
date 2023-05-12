package handler

type AddPositionRequest struct {
	Position string `json:"position" validate:"required,min=3"`
	Tag      string `json:"tag" validate:"required,min=2"`
}
