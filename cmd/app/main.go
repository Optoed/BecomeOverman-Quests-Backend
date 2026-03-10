package main

import (
	"BecomeOverMan/internal/handlers"
	"BecomeOverMan/internal/kafka"
	"log"
	"log/slog"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"BecomeOverMan/internal/config"
	_ "BecomeOverMan/internal/models"
	"BecomeOverMan/internal/repositories"
	"BecomeOverMan/internal/services"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	db, err := sqlx.Connect("postgres", config.Cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	techRepo := repositories.NewTechRepository(db)
	techService := services.NewTechService(techRepo)

	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)

	questRepo := repositories.NewQuestRepository(db)
	kafkaProducer := kafka.NewProducerFromEnv()
	if kafkaProducer == nil {
		slog.Info("Kafka producer disabled: KAFKA_BROKERS is empty")
	} else {
		defer func() {
			if err := kafkaProducer.Close(); err != nil {
				slog.Error("Failed to close kafka producer", "error", err)
			}
		}()
	}

	questService := services.NewQuestService(questRepo, userRepo, kafkaProducer)

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service": "BecomeOverMan Quests Backend",
			"status":  "ok",
			"routes": []string{
				"/health/db",
				"/auth/register",
				"/auth/login",
				"/quests/shop",
				"/users/me/quests/:questID",
			},
		})
	})

	handlers.RegisterTechRoutes(r, techService)
	handlers.RegisterAuthRoutes(r, userService)
	handlers.RegisterUserRoutes(r, userService)
	handlers.RegisterQuestRoutes(r, questService)

	if err := r.Run("0.0.0.0:8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
