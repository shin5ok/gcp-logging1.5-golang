package main

import (
	"context"
	"log"
	"math/rand"
	"os"
	"time"

	"net/http"

	logging "cloud.google.com/go/logging"
	"github.com/gin-gonic/gin"
)

var (
	appName = "tiny-api"
)

func main() {

	ctx := context.Background()

	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	client, err := logging.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	logger := client.Logger(appName, logging.RedirectAsJSON(os.Stderr))

	g := gin.Default()

	g.GET("/api", func(c *gin.Context) {
		start := time.Now()
		rand.Seed(time.Now().UnixNano())
		n := rand.Intn(100)
		time.Sleep(time.Duration(n) * time.Millisecond)
		finish := time.Now()
		duration := finish.Sub(start)
		logger.Log(
			logging.Entry{
				Severity: logging.Info,
				Payload: map[string]interface{}{
					"Duration": duration.Seconds(),
					"RemoteIP": c.RemoteIP(),
					"Service":  os.Getenv("K_REVISION"),
				},
			})
		c.JSON(http.StatusOK, gin.H{})
	})

	g.Run(":8080")
}
