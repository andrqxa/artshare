package router

import (
	usercontroller "artshare/internal/controller/user"
	userrouter "artshare/internal/router/user"

	"github.com/go-chi/chi/v5"
)

type Router struct {
	r           *chi.Mux
	userHandler *userrouter.Handler
}

func NewRouter() *Router {
	r := chi.NewRouter()
	return &Router{
		r:           r,
		userHandler: userrouter.NewHandler(usercontroller.NewController()),
	}
}

func (rr *Router) Mount() {
	rr.r.Route("/api/v1", func(r chi.Router) {
		r.Get("/users/me", rr.userHandler.GetCurrentUser)
		r.Patch("/users/me", rr.userHandler.UpdateCurrentUser)
	})
}

func (rr *Router) Router() *chi.Mux {
	return rr.r
}
