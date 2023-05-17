package handler

type GetAllPositionResponse struct {
	PositionID int    `json:"position_id"`
	Position   string `json:"position"`
	Tag        string `json:"tag"`
}

type Meta struct {
	CurrentLimit  int `json:"current_limit"`
	CurrentOffset int `json:"current_offset"`
	CurrentPage   int `json:"current_page"`
	TotalData     int `json:"total_data"`
	TotalPage     int `json:"total_page"`
}
