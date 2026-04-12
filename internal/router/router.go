package router

import (
	"artshare/internal/router/user"

	"github.com/go-chi/chi/v5"
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
	rr.r.Post("/users/me", rr.userHandler.CreateUser)
	rr.r.Patch("/users.me", rr.userHandler.UpdateUser)
}

func (rr *Router) Router() *chi.Mux {
	return rr.r
}
