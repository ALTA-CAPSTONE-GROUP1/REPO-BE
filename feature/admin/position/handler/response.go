package handler

type GetAllPositionResponse struct {
	PositionID int    `json:"position_id"`
	Position    string `json:"position"`
	Tag         string `json:"tag"`
}
