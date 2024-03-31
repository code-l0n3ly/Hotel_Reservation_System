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
)

// App encapsulates Environment, Router, and DB connections
type App struct {
	Router        *mux.Router
	DB            *Database.DBExecutor
	UserHandler   *Handlers.UserHandler
	ReviewHandler *Handlers.ReviewHandler
	UnitHandler   *Handlers.UnitHandler
}

// Initialize sets up the database connection and the router
func (a *App) Initialize(user, password, dbname string) {
	connectionString := fmt.Sprintf("%s:%s@/%s", user, password, dbname)
	var err error
	a.DB, err = Database.InitDB(connectionString)
	if err != nil {
		log.Fatal(err)
	}
	a.Router = mux.NewRouter()
	a.UserHandler = Handlers.NewUserHandler(a.DB.Db)
	a.UnitHandler = Handlers.NewUnitHandler(a.DB.Db)
	a.ReviewHandler = Handlers.NewReviewHandler(a.DB.Db)
	a.initializeRoutes()
}

// InitializeRoutes sets up the routes for the application
func (a *App) initializeRoutes() {
	Routes.RegisterUserRoutes(a.Router, a.UserHandler)
	Routes.RegisterReviewRoutes(a.Router, a.ReviewHandler)
	Routes.RegisterUnitRoutes(a.Router, a.UnitHandler)
}

// Run starts the server on a specified port
func (a *App) Run(addr string) {
	log.Printf("Listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

// In app.go
