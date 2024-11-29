package main

import (
	"log"
	httpNet "net/http"

	"blog-api/config"
	"blog-api/internal/delivery/http"
	"blog-api/internal/repository/mysql"
	"blog-api/internal/usecase"
	"blog-api/pkg/db"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	// Load config
	cfg := config.LoadConfig()
	dbConn := db.ConnectDB(cfg)
	if err := db.InitializeDB(dbConn, cfg); err != nil {
		log.Fatalf("Error initializing the database: %v", err)
	}

	userRepo := mysql.NewUserRepository(dbConn)
	userUsecase := usecase.NewUserUsecase(userRepo, cfg.JWTSecret)
	log.Println(userUsecase)

	blogRepo := mysql.NewBlogRepository(dbConn)
	blogUsecase := usecase.NewBlogUsecase(blogRepo)

	r := mux.NewRouter()

	http.NewUserHandler(r, userUsecase)
	http.NewBlogHandler(r, blogUsecase, config.LoadConfig().JWTSecret)

	log.Println("Server is running on port 8080")
	log.Fatal(httpNet.ListenAndServe(":8080", r))
}
