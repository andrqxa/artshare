package router

import (
	"os/user"

	"github.com/go-chi/chi"
)

type Router struct {
	r           *chi.Mux
	userHandler *user.Handler
}

func NewRouter() *Router {
	r := chi.NewRouter()
	return &Router{r: r}
}

func (rr *Router) Mount() {
	rr.Post("/users", rr.userHandler.CreateUser)
}

func (rr *Router) Router() *chi.Mux {
	return rr.r
}
