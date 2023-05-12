package handler

type AddPositionRequest struct {
	Position string `json:"position"`
	Tag      string `json:"tag"`
}
