package handler

import "mime/multipart"

type AddAddSubReq struct {
	To              []string        `form:"to" json:"to" xml:"to" yaml:"to" query:"to"`
	CC              []string        `form:"cc" json:"cc" xml:"cc" yaml:"cc" query:"cc"`
	SubmissionType  string          `form:"submission_type" json:"submission_type" xml:"submission_type" yaml:"submission_type" query:"submission_type"`
	SubmissionValue int             `form:"submission_value" json:"submission_value" xml:"submission_value" yaml:"submission_value" query:"submission_value"`
	Title           string          `form:"title" json:"title" xml:"title" yaml:"title" query:"title"`
	Message         string          `form:"message" json:"message" xml:"message" yaml:"message" query:"message"`
	Attachment      *multipart.File `form:"attachment" json:"attachment" xml:"attachment" yaml:"attachment" query:"attachment"`
}


type UpdateSubReq struct {
	To              []string        `form:"to" json:"to" xml:"to" yaml:"to" query:"to"`
	CC              []string        `form:"cc" json:"cc" xml:"cc" yaml:"cc" query:"cc"`
	
	SubmissionType  string          `form:"submission_type" json:"submission_type" xml:"submission_type" yaml:"submission_type" query:"submission_type"`
	SubmissionValue int             `form:"submission_value" json:"submission_value" xml:"submission_value" yaml:"submission_value" query:"submission_value"`
	Title           string          `form:"title" json:"title" xml:"title" yaml:"title" query:"title"`
	Message         string          `form:"message" json:"message" xml:"message" yaml:"message" query:"message"`
	Attachment      *multipart.File `form:"attachment" json:"attachment" xml:"attachment" yaml:"attachment" query:"attachment"`
}