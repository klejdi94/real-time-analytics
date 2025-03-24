package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/user/analytic-dashboard/pkg/data"
	"github.com/user/analytic-dashboard/pkg/processing"
	"github.com/user/analytic-dashboard/pkg/visualization"
)

func main() {
	// Set up router
	r := gin.Default()

	// Initialize components
	dataService := data.NewService()
	processor := processing.NewProcessor(dataService)
	vizManager := visualization.NewManager()

	// Connect the processor to the visualization manager
	processor.AddListener(vizManager.BroadcastMetrics)

	// Start the visualization manager in a goroutine
	go vizManager.Start()

	// Set up static files
	r.Static("/static", "./web/static")
	r.LoadHTMLGlob("web/templates/*")

	// Set up routes
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Real-time Analytics Dashboard",
		})
	})

	// API routes
	api := r.Group("/api")
	{
		// Data ingestion endpoint
		api.POST("/ingest", func(c *gin.Context) {
			var payload data.Payload
			if err := c.ShouldBindJSON(&payload); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			if err := dataService.Store(payload); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"status": "data received"})
		})

		// Data retrieval endpoints
		api.GET("/metrics", func(c *gin.Context) {
			metrics, err := processor.GetMetrics()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, metrics)
		})

		// WebSocket for real-time updates
		api.GET("/ws", func(c *gin.Context) {
			vizManager.HandleWebSocket(c.Writer, c.Request)
		})
	}

	// Start server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Start the server in a goroutine
	go func() {
		log.Println("Server starting on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Start the data processor in a goroutine
	go processor.Start()

	// Set up graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Give the server 5 seconds to finish ongoing requests
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	// Clean up resources
	processor.Stop()
	log.Println("Server exited properly")
}
