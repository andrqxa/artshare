package router

import (
	artistcontroller "github.com/andrqxa/artshare/internal/controller/artist"
	artworkcontroller "github.com/andrqxa/artshare/internal/controller/artwork"
	exchangecontroller "github.com/andrqxa/artshare/internal/controller/exchange"
	likecontroller "github.com/andrqxa/artshare/internal/controller/like"
	usercontroller "github.com/andrqxa/artshare/internal/controller/user"
	artistrouter "github.com/andrqxa/artshare/internal/router/artist"
	artworkrouter "github.com/andrqxa/artshare/internal/router/artwork"
	exchangerouter "github.com/andrqxa/artshare/internal/router/exchange"
	likerouter "github.com/andrqxa/artshare/internal/router/like"
	userrouter "github.com/andrqxa/artshare/internal/router/user"

	"github.com/go-chi/chi/v5"
)

type Router struct {
	r               *chi.Mux
	userHandler     *userrouter.Handler
	artistHandler   *artistrouter.Handler
	artworkHandler  *artworkrouter.Handler
	likeHandler     *likerouter.Handler
	exchangeHandler *exchangerouter.Handler
}

func NewRouter() *Router {
	r := chi.NewRouter()
	return &Router{
		r:               r,
		userHandler:     userrouter.NewHandler(usercontroller.NewController()),
		artistHandler:   artistrouter.NewHandler(artistcontroller.NewController()),
		artworkHandler:  artworkrouter.NewHandler(artworkcontroller.NewController()),
		likeHandler:     likerouter.NewHandler(likecontroller.NewController()),
		exchangeHandler: exchangerouter.NewHandler(exchangecontroller.NewController()),
	}
}

func (rr *Router) Mount() {
	rr.r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/register", rr.userHandler.RegisterUser)
		r.Post("/auth/login", rr.userHandler.LoginUser)

		r.Get("/users/me", rr.userHandler.GetCurrentUser)
		r.Patch("/users/me", rr.userHandler.UpdateCurrentUser)

		r.Post("/artists/me", rr.artistHandler.CreateCurrentArtist)
		r.Get("/artists/me", rr.artistHandler.GetCurrentArtist)
		r.Patch("/artists/me", rr.artistHandler.UpdateCurrentArtist)
		r.Get("/artists", rr.artistHandler.ListArtists)
		r.Get("/artists/{artistId}", rr.artistHandler.GetArtistByID)

		r.Get("/artworks", rr.artworkHandler.ListArtworks)
		r.Post("/artworks", rr.artworkHandler.CreateArtwork)
		r.Get("/artworks/{artworkId}", rr.artworkHandler.GetArtworkByID)
		r.Patch("/artworks/{artworkId}", rr.artworkHandler.UpdateArtwork)
		r.Delete("/artworks/{artworkId}", rr.artworkHandler.DeleteArtwork)

		r.Post("/artworks/{artworkId}/likes", rr.likeHandler.LikeArtwork)
		r.Delete("/artworks/{artworkId}/likes", rr.likeHandler.UnlikeArtwork)

		r.Get("/exchange-requests", rr.exchangeHandler.ListExchangeRequests)
		r.Post("/exchange-requests", rr.exchangeHandler.CreateExchangeRequest)
		r.Get("/exchange-requests/{exchangeRequestId}", rr.exchangeHandler.GetExchangeRequestByID)
		r.Patch("/exchange-requests/{exchangeRequestId}", rr.exchangeHandler.UpdateExchangeRequestStatus)
	})
}

func (rr *Router) Router() *chi.Mux {
	return rr.r
}
