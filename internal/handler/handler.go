package handler

type Handler struct {
	Auth *Auth
}

func New() *Handler {
	return &Handler{}
}
