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

	g.GET("/", func(c *gin.Context) {
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
					"duration":         duration.Seconds(),
					"client_ip":        c.Request.Header.Get("X-Forwarded-For"),
					"service_revision": os.Getenv("K_REVISION"),
					"accept":           c.Request.Header.Get("Accept"),
					"method":           c.Request.Method,
					"host":             c.Request.Host,
					"misc":             "test",
					"path":             c.Request.URL.Path,
				},
			})
		c.JSON(http.StatusOK, gin.H{})
	})

	type LogStruct struct {
		Duration        string `json:"duration"`
		ClientIp        string `json:"client_ip"`
		ServiceRevision string `json:"service_revision"`
		Accept          string `json:"accept"`
		Method          string `json:"method"`
		Host            string `json:"host"`
		Misc            string `json:"misc"`
		Path            string `json:"path"`
	}
	g.GET("/test", func(c *gin.Context) {
		logger.Log(
			logging.Entry{
				Severity: logging.Warning,
				Payload: LogStruct{
					Duration:        "",
					ClientIp:        c.Request.Header.Get("X-Forwarded-For"),
					ServiceRevision: os.Getenv("K_REVISION"),
					Accept:          c.Request.Header.Get("Accept"),
					Method:          c.Request.Method,
					Host:            c.Request.Host,
					Misc:            "test",
					Path:            c.Request.URL.Path,
				},
			})
		c.JSON(http.StatusOK, gin.H{})
	})

	g.Run(":8080")
}
