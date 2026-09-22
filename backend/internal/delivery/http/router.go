package httpapi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	"github.com/cassianobraz/SistemaControleChamados/backend/internal/delivery/http/handler"
	"github.com/cassianobraz/SistemaControleChamados/backend/internal/delivery/http/middleware"
	"github.com/cassianobraz/SistemaControleChamados/backend/internal/usecase"
)

type Dependencies struct {
	TicketUseCase      usecase.TicketUseCase
	ResponsibleUseCase usecase.ResponsibleUseCase
	OpenAPISpec        []byte
}

func NewRouter(deps Dependencies) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Logging, middleware.CORS)

	ticketHandler := handler.NewTicketHandler(deps.TicketUseCase)
	responsibleHandler := handler.NewResponsibleHandler(deps.ResponsibleUseCase)

	router.Get("/health", handler.Health)

	router.Route("/api/v1/tickets", func(r chi.Router) {
		r.Post("/", ticketHandler.Create)
		r.Get("/", ticketHandler.List)
		r.Get("/{id}", ticketHandler.Get)
		r.Put("/{id}", ticketHandler.Update)
		r.Patch("/{id}/assign", ticketHandler.Assign)
	})

	router.Route("/api/v1/responsibles", func(r chi.Router) {
		r.Get("/", responsibleHandler.List)
		r.Post("/", responsibleHandler.Create)
		r.Put("/{id}", responsibleHandler.Update)
		r.Delete("/{id}", responsibleHandler.Delete)
	})

	router.Get("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write(deps.OpenAPISpec)
	})
	router.Handle("/swagger/*", httpSwagger.Handler(httpSwagger.URL("/swagger/doc.json")))

	return router
}
