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
	From           Sender
	To             Receiver
	Title          string
	SubmissionType string
	Attachment     string
}

type Sender struct {
	Name     string
	Position string
}

type Receiver struct {
	Name     string
	Position string
}
