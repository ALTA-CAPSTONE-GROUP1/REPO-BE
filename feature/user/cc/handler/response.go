package handler

type Response struct {
	SubmissionID   int      `json:"submission_id"`
	From           Sender   `json:"from"`
	To             Receiver `json:"to"`
	Title          string   `json:"title"`
	SubmissionType string   `json:"submission_type"`
	Attachment     string   `json:"attachment"`
}

type Sender struct {
	Name     string `json:"name"`
	Position string `json:"position"`
}

type Receiver struct {
	Name     string `json:"name"`
	Position string `json:"position"`
}
