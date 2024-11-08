package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aminGhafoory/daq/controllers"
	"github.com/aminGhafoory/daq/internal/database"
	"github.com/aminGhafoory/daq/models"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/csrf"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
)

type Config struct {
	listenAddr string
	DBURL      string
	CSRF       struct {
		Key    string
		Secure bool
	}
}

func LoadEnvConfig() (Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return Config{}, err
	}

	listenAddr := os.Getenv("LISTEN_ADDR")
	if listenAddr == "" {
		return Config{}, err
	}

	dbURL := os.Getenv("DB_URL")
	if listenAddr == "" {
		return Config{}, err
	}

	CSRFKey := os.Getenv("CSRF")
	if listenAddr == "" {
		return Config{}, err
	}

	devEnv := os.Getenv("DEV")
	if devEnv == "" {
		return Config{}, fmt.Errorf("CSRFKey Not found in the .env file")

	}
	useSecureMode := func(isDev string) bool {
		if isDev == "false" {
			return true
		} else {
			return false
		}
	}(devEnv)

	c := Config{}
	c.listenAddr = listenAddr
	c.CSRF.Key = CSRFKey
	c.CSRF.Secure = useSecureMode
	c.DBURL = dbURL

	return c, nil

}

func main() {
	fmt.Println("hello")

	config, err := LoadEnvConfig()
	if err != nil {
		log.Fatalf("Could Not Read .env file %+v", err)
	}

	//MiddleWare
	CSRFMiddleware := csrf.Protect(
		[]byte(config.CSRF.Key),
		csrf.Secure(config.CSRF.Secure),
		csrf.CookieName("CSRF"),
		csrf.Path("/"),
		csrf.HttpOnly(true),
	)

	//DB
	conn, err := sql.Open("pgx", config.DBURL)
	if err != nil {
		log.Fatal("Can NOT connect to the database")
	}
	db := database.New(conn)
	err = conn.Ping()
	if err != nil {
		log.Fatal("can not connect to DB")
	}

	//CONTROLLERS

	UserC := controllers.Users{
		UserService: &models.UserService{
			DB: db,
		},
		SessionService: &models.SessionService{
			DB: db,
		},
	}

	UserMw := controllers.UserMiddleware{
		SessionService: &models.SessionService{
			DB: db,
		},
	}

	r := chi.NewMux()

	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)
	r.Use(CSRFMiddleware)
	r.Use(UserMw.SetUser)

	r.Use(CSRFMiddleware)

	// r.Get("/sign-in", func(w http.ResponseWriter, r *http.Request) {

	// 	component := signIn.SignIn("hello", r, []string{})
	// 	component.Render(context.Background(), w)
	// })

	// r.Get("/sign-up", func(w http.ResponseWriter, r *http.Request) {

	// 	component := signUp.NewUser("hello", r, []string{})
	// 	component.Render(context.Background(), w)
	// })

	r.Get("/public/*", func(w http.ResponseWriter, r *http.Request) {
		fs := http.FileServer(http.Dir("."))
		fs.ServeHTTP(w, r)
	})

	r.Post("/sign-in/", UserC.ProccessSignIn)
	r.Post("/users", UserC.Create)
	r.Get("/sign-in", UserC.NewSignInPage)
	r.Get("/sign-out", UserC.SignOutUser)
	r.Post("/sign-in", UserC.ProccessSignIn)
	r.Get("/sign-up", UserC.New)

	fmt.Printf("server started on http://localhost%s\n", config.listenAddr)
	err = http.ListenAndServe(config.listenAddr, r)
	if err != nil {
		log.Fatal(err)
	}

}
