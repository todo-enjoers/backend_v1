package model

type Todo struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsDone      bool   `json:"isDone"`
}

type TodoCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsDone      bool   `json:"isDone"`
}

type UserDTO struct {
	ID       int64  `json:"id"`
	Login    string `json:"email"`
	Password string `json:"password"`
}
