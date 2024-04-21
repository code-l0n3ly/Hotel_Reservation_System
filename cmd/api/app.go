// app.go
package App

import (
	"fmt"
	"log"
	"net/http"

	Routes "GraduationProject.com/m/internal/Routes"
	Database "GraduationProject.com/m/internal/db"
	Handlers "GraduationProject.com/m/internal/handler"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// App encapsulates Environment, Router, and DB connections
type App struct {
	Router                      *mux.Router
	DB                          *Database.DBExecutor
	UserHandler                 *Handlers.UserHandler
	ReviewHandler               *Handlers.ReviewHandler
	UnitHandler                 *Handlers.UnitHandler
	BookingHandler              *Handlers.BookingHandler
	ReportHandler               *Handlers.ReportHandler
	FinancialTransactionHandler *Handlers.FinancialTransactionHandler
	MaintenanceTicketHandler    *Handlers.MaintenanceTicketHandler
	PropertyHandler             *Handlers.PropertyHandler
	MessageHandler              *Handlers.MessageHandler
}

// Initialize sets up the database connection and the router
func (a *App) Initialize(user, password, address, dbname string) {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s", user, password, address, dbname)
	var err error
	a.DB, err = Database.InitDB(connectionString)
	if err != nil {
		log.Fatal(err)
	}
	a.Router = mux.NewRouter()
	a.UserHandler = Handlers.NewUserHandler(a.DB.Db)
	a.UnitHandler = Handlers.NewUnitHandler(a.DB.Db)
	a.ReviewHandler = Handlers.NewReviewHandler(a.DB.Db)
	a.BookingHandler = Handlers.NewBookingHandler(a.DB.Db)
	a.ReportHandler = Handlers.NewReportHandler(a.DB.Db)
	a.FinancialTransactionHandler = Handlers.NewFinancialTransactionHandler(a.DB.Db)
	a.MaintenanceTicketHandler = Handlers.NewMaintenanceTicketHandler(a.DB.Db)
	a.PropertyHandler = Handlers.NewPropertyHandler(a.DB.Db)
	a.MessageHandler = Handlers.NewMessageHandler(a.DB.Db)
	a.initializeRoutes()
}

// InitializeRoutes sets up the routes for the application
func (a *App) initializeRoutes() {
	Routes.RegisterUserRoutes(a.Router, a.UserHandler)
	Routes.RegisterReviewRoutes(a.Router, a.ReviewHandler)
	Routes.RegisterUnitRoutes(a.Router, a.UnitHandler)
	Routes.RegisterBookingRoutes(a.Router, a.BookingHandler)
	Routes.RegisterReportRoutes(a.Router, a.ReportHandler)
	Routes.RegisterFinancialTransactionRoutes(a.Router, a.FinancialTransactionHandler)
	Routes.RegisterMaintenanceTicketRoutes(a.Router, a.MaintenanceTicketHandler)
	Routes.RegisterPropertyRoutes(a.Router, a.PropertyHandler)
	Routes.RegisterMessageRoutes(a.Router, a.MessageHandler)

}

// Run starts the server on a specified port
func (a *App) Run(addr string) {
	log.Printf("Listening on %s\n", addr)
	handler := cors.AllowAll().Handler(a.Router)
	log.Fatal(http.ListenAndServe(addr, handler))
}

// In app.go
