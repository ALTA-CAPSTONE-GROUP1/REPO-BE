package cc

import "github.com/labstack/echo/v4"

type Handler interface {
	GetAllCcHander() echo.HandlerFunc
}

type UseCase interface {
	GetAllCcLogic(userID string) ([]CcCore, error)
}

type Repository interface {
	GetAllCc(userID string) ([]CcCore, error)
}

type CcCore struct {
	SubmisisonID   int
	Title          string
	SubmissionType string
	Attachment     string
	From           Sender
	To             Receiver
}

type Sender struct {
	Name     string
	Position string
}

type Receiver struct {
	Name     string
	Position string
}

type QueryParams struct {
	Limit  int
	Offset int
}
