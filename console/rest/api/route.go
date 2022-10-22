package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Run run http router
func (a *API) Run(port string) *http.Server {
	// init gin
	gin.SetMode(gin.DebugMode)
	router := gin.New()

	router.GET("/healthz", a.handleHealthCheck)

	v1 := router.Group("/v1")
	{
		v1.GET("/index", a.handleIndexerCommand)
		v1.GET("/block", a.handleGetLatestBlock)
		v1.GET("/block/:id", a.handleGetBlock)
		v1.GET("/stats", a.handleGetStats)
		v1.GET("/stats/:range", a.handleGetRangeStats)
		v1.GET("/tx", a.handleGetLatestTx)
		v1.GET("/tx/:hash", a.handleGetTx)
	}

	// gin middleware config
	router.Use(gin.Recovery())
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "HEAD", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "X-Requested-With", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"status": false, "data": gin.H{"status": false, "message": fmt.Sprintf("Page not found: %s, method: %s", c.Request.URL, c.Request.Method)}})
	})

	router.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"status": false, "data": gin.H{"status": false, "message": "Method not found"}})
	})

	// Note: we use http server to have graceful shutdown
	server := &http.Server{
		Addr:         port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 20 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	go func() {
		log.Printf("Listening and serving HTTP on %s\n", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("gin sever stoped with err: %s \n", err)
		}
	}()

	return server
}

// Response is the struct we return the response
type Response struct {
	Status  int         `json:"status"`
	Payload interface{} `json:"payload"`
}

func (a *API) response(c *gin.Context, status int, payload interface{}) {
	response := Response{
		Status:  status,
		Payload: payload,
	}

	c.Header("Content-Type", "application/json")
	c.Status(status)

	c.JSON(status, response)
}
