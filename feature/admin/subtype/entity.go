package subtype

import (
	"github.com/labstack/echo/v4"
)

type Core struct {
	SubmissionTypeName string
	PositionTag        []string
	SubmissionValues   []ValueDetails
	Requirement        string
}

type ValueDetails struct {
	Value         int
	TagPositionTo []string
	TagPositionCC []string
}

type RepoData struct {
	TypeName               string
	TypeRequirement        string
	OwnersTag              []string
	SubTypeInterdependence []RepoDataInterdependence
}

type RepoDataInterdependence struct {
	Value  int
	TosTag []string
	CcsTag []string
}

type Handler interface {
	AddTypeHandler() echo.HandlerFunc
}

type UseCase interface {
	AddSubTypeLogic(newType Core) error
}

type Repository interface {
	InsertSubType(req RepoData) error
}
