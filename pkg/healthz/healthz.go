package healthz

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kaijun123/kubernetes-kms/pkg/util"
)

func InitHttpServer(service util.Service) *http.Server {
	httpServer := gin.Default()
	httpServer.GET("/healthz", func(c *gin.Context) {
		log.Println("Start of Health Check......")
		status, err := service.Status(context.Background())

		if err != nil || status.Healthz != "ok" {
			err := fmt.Errorf("failed health check: %s", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"healthcheck": "ok",
		})
		log.Println("End of Health Check......")
	})

	// Return the http.Server instance to be managed by the main function
	return &http.Server{
		Addr:    ":8087",
		Handler: httpServer,
	}
}
