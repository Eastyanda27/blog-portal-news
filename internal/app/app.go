package app

import (
	"blognewsportal/config"
	"blognewsportal/internal/adapter/handler"
	"blognewsportal/internal/adapter/repository"
	"blognewsportal/internal/core/service"
	"blognewsportal/lib/auth"
	"blognewsportal/lib/middleware"
	"blognewsportal/lib/pagination"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func RunServer() {
	cfg := config.NewConfig()
	fmt.Printf("CONFIG LOADED: %+v\n", cfg.Psql)
	db, err := cfg.ConnectionPostgres()
	if err != nil {
		log.Fatalf("Error to connect to database: %v", err)
		return
	}

	err = os.MkdirAll("./temp/content", 0755)
	if err != nil {
		log.Fatalf("Error creating temp directory: %v", err)
	}

	jwt := auth.NewJwt(cfg)

	middlewareAuth := middleware.NewMiddleware(cfg)

	_ = pagination.NewPagination()

	authRepo := repository.NewAuthRepository(db.DB)
	categoryRepo := repository.NewCategoryRepository(db.DB)
	contentRepo := repository.NewContentRepository(db.DB)
	userRepo := repository.NewUserRepository(db.DB)

	authService := service.NewAuthService(authRepo, cfg, jwt)
	categoryService := service.NewCategoryService(categoryRepo)
	contentService := service.NewContentService(contentRepo)
	userService := service.NewUserService(userRepo)

	authHandler := handler.NewAuthHandler(authService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	contentHandler := handler.NewContentHandler(contentService)
	useHandler := handler.NewUserHandler(userService)

	app := fiber.New()
	app.Use(cors.New())
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${ip} ${status} - ${latency} ${method} ${path}\n",
	}))

	if os.Getenv("APP_ENV") != "production" {
		cfg := swagger.Config{
			BasePath: "/api",
			FilePath: "./docs/swagger.json",
			Path:     "docs",
			Title:    "Swagger API Docs",
		}

		app.Use(swagger.New(cfg))
	}

	api := app.Group("/api")

	api.Post("/login", authHandler.Login)

	adminApp := api.Group("/admin")
	adminApp.Use(middlewareAuth.CheckToken())

	categoryApp := adminApp.Group("/categories")
	categoryApp.Get("/", categoryHandler.GetCategories)
	categoryApp.Post("/", categoryHandler.CreateCategory)
	categoryApp.Get("/:categoryID", categoryHandler.GetCategoryByID)
	categoryApp.Put("/:categoryID", categoryHandler.EditCategoryByID)
	categoryApp.Delete("/:categoryID", categoryHandler.DeleteCategory)

	contentApp := adminApp.Group("/contents")
	contentApp.Get("/", contentHandler.GetContents)
	contentApp.Post("/", contentHandler.CreateContent)
	contentApp.Get("/:contentID", contentHandler.GetContentByID)
	contentApp.Put("/:contentID", contentHandler.EditContentByID)
	contentApp.Delete("/:contentID", contentHandler.DeleteContent)

	userApp := adminApp.Group("/users")
	userApp.Get("/profile", useHandler.GetUserByID)
	userApp.Put("/update-password", useHandler.EditPassword)

	feApp := api.Group("/fe")
	feApp.Get("/categories", categoryHandler.GetCategoryFE)
	feApp.Get("/contents", contentHandler.GetContentWithQuery)
	feApp.Get("/contents/:contentID", contentHandler.GetContentDetail)

	go func() {
		if cfg.App.AppPort == "" {
			cfg.App.AppPort = os.Getenv("APP_PORT")
		}

		err := app.Listen(":" + cfg.App.AppPort)
		if err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGTERM)

	<-quit

	log.Println("server shutdown of 5 seconds")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	app.ShutdownWithContext(ctx)
}
