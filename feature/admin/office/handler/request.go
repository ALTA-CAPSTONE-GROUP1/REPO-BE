package handler

type AddOfficeRequest struct {
	Name     string `json:"name"`
	Level    string `json:"level"`
	ParentID uint   `json:"parent_id"`
}
