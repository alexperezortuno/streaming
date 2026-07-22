package handler

import (
	"github.com/alexperezortuno/streaming/internal/config"
	"github.com/alexperezortuno/streaming/internal/media"
	"github.com/alexperezortuno/streaming/internal/middleware"
	"github.com/alexperezortuno/streaming/internal/repository"
	"github.com/alexperezortuno/streaming/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handlers struct {
	Auth   *AuthHandler
	Video  *VideoHandler
	Stream *StreamHandler
	List   *ListHandler
	Health *HealthHandler
}

func NewHandlers(pool *pgxpool.Pool, cfg *config.Config) *Handlers {
	userRepo := repository.NewUserRepository(pool)
	videoRepo := repository.NewVideoRepository(pool)
	listRepo := repository.NewListRepository(pool)

	authService := service.NewAuthService(userRepo, cfg)
	listService := service.NewListService(listRepo)
	streamService := service.NewStreamService(cfg)

	transcoder := media.NewTranscoder(cfg.TranscodeWorkers)
	videoService := service.NewVideoService(videoRepo, cfg, transcoder)

	return &Handlers{
		Auth:   NewAuthHandler(authService),
		Video:  NewVideoHandler(videoService),
		Stream: NewStreamHandler(streamService),
		List:   NewListHandler(listService),
		Health: NewHealthHandler(),
	}
}

func NewRouter(h *Handlers, cfg *config.Config, pool *pgxpool.Pool) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Recovery)
	r.Use(middleware.Logger)
	r.Use(middleware.CORS(cfg.CORSOrigins))

	userRepo := repository.NewUserRepository(pool)
	authService := service.NewAuthService(userRepo, cfg)

	r.Get("/api/health", h.Health.Check)

	r.Route("/api/auth", func(r chi.Router) {
		r.Post("/login", h.Auth.Login)
		r.Post("/register", h.Auth.Register)
	})

	r.Route("/api/lists", func(r chi.Router) {
		r.Use(middleware.Authenticator(authService))
		r.Get("/", h.List.List)
		r.Post("/", h.List.Create)
		r.Get("/{id}", h.List.Get)
		r.Delete("/{id}", h.List.Delete)
	})

	r.Route("/api/videos", func(r chi.Router) {
		r.Use(middleware.Authenticator(authService))
		r.Get("/", h.Video.List)
		r.Post("/", h.Video.Upload)
		r.Get("/{id}", h.Video.Get)
		r.Put("/{id}", h.Video.Update)
		r.Delete("/{id}", h.Video.Delete)
		r.Get("/{id}/stream", h.Stream.Serve)
		r.Get("/{id}/stream/{seg}", h.Stream.Serve)
	})

	return r
}
