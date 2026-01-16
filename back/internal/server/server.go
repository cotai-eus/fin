package server

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
	"github.com/lauratech/fin/back/internal/config"
	"github.com/lauratech/fin/back/internal/modules/bills"
	"github.com/lauratech/fin/back/internal/modules/budgets"
	"github.com/lauratech/fin/back/internal/modules/cards"
	"github.com/lauratech/fin/back/internal/modules/support"
	"github.com/lauratech/fin/back/internal/modules/transfers"
	"github.com/lauratech/fin/back/internal/modules/users"
)

// Server holds dependencies for HTTP server
type Server struct {
	Config           *config.Config
	DB               *sql.DB
	router           *chi.Mux
	usersHandler     *users.Handler
	transfersHandler *transfers.Handler
	cardsHandler     *cards.Handler
	billsHandler     *bills.Handler
	budgetsHandler   *budgets.Handler
	supportHandler   *support.Handler
}

// New creates a new server instance
func New(cfg *config.Config, db *sql.DB) *Server {
	// Initialize repositories
	usersRepo := users.NewRepository(db)
	transfersRepo := transfers.NewRepository(db)
	cardsRepo := cards.NewRepository(db, cfg.EncryptionKey)
	billsRepo := bills.NewRepository(db)
	budgetsRepo := budgets.NewRepository(db)
	supportRepo := support.NewRepository(db)

	// Initialize services
	usersService := users.NewService(usersRepo)
	transfersService := transfers.NewService(transfersRepo, usersRepo, db)
	cardsService := cards.NewService(cardsRepo, db)
	billsService := bills.NewService(billsRepo, db)
	budgetsService := budgets.NewService(budgetsRepo, db)
	supportService := support.NewService(supportRepo, db)

	// Initialize handlers
	usersHandler := users.NewHandler(usersService)
	transfersHandler := transfers.NewHandler(transfersService)
	cardsHandler := cards.NewHandler(cardsService)
	billsHandler := bills.NewHandler(billsService)
	budgetsHandler := budgets.NewHandler(budgetsService)
	supportHandler := support.NewHandler(supportService)

	s := &Server{
		Config:           cfg,
		DB:               db,
		usersHandler:     usersHandler,
		transfersHandler: transfersHandler,
		cardsHandler:     cardsHandler,
		billsHandler:     billsHandler,
		budgetsHandler:   budgetsHandler,
		supportHandler:   supportHandler,
	}

	s.router = s.setupRouter()

	return s
}

// Router returns the configured Chi router
func (s *Server) Router() *chi.Mux {
	return s.router
}
