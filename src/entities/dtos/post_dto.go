package dtos

type PostDTO struct {
	Message string `json:"message"`
	Sender  string `json:"sender"`
	Room    string `json:"room"`
}
