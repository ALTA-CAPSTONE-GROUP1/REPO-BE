package handler

type GetAllPositionResponse struct {
	PositionID int    `json:"position_id"`
	Position   string `json:"position"`
	Tag        string `json:"tag"`
}

type Meta struct {
	Current_limit int `json:"current_limit"`
	Current_Page  int `json:"current_page"`
	Data_amount   int `json:"data_amount"`
	Page_amount   int `json:"page_amount"`
}
