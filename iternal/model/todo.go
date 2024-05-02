package model

type TodoDTO struct {
	ID          int
	Name        string
	Description string
	IsDone      bool
}

type TodoCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
