package main

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

func init() {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		slog.Error("Reading env from .env file failed, using os environment variables", "error", err)
		viper.AutomaticEnv()
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoURI := viper.GetString("MONGODB_URI")

	var err error
	mongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		slog.Error("Failed to connect to MongoDB", "error", err, "uri", mongoURI)
		return
	}

	r := gin.Default()

	r.GET("/healthz", healthCheckHandler)

	port := viper.GetString("APP_PORT")
	if err := r.Run(":" + port); err != nil {
		slog.Error("Failed to start server", "error", err, "port", port)
	} else {
		slog.Info("Server is running on port", "port", port)
	}
}

// Handler untuk /healthz
func healthCheckHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := mongoClient.Ping(ctx, nil)
	if err != nil {
		slog.Error("MongoDB is not healthy", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  err.Error(),
			"status": "unhealthy",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
