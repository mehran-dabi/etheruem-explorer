package api

import (
	"energi-challenge/config"
	"energi-challenge/console/indexer/utils"
	"energi-challenge/infrastructure/sqlite"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Router struct {
	store       sqlite.ISqlite
	jobs        chan int64
	conf        config.Configs
	latestBlock int64
}

// Run runs http router
func (idx *Indexer) run(port string) *http.Server {
	// init gin
	gin.SetMode(gin.DebugMode)
	router := gin.New()

	router.GET("/health", idx.handleHealthCheck)

	router.GET("/", idx.handleIndex)

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

// handleHealthCheck checks sqlite ping
func (idx *Indexer) handleHealthCheck(c *gin.Context) {
	health := map[string]interface{}{
		"store": "up",
	}

	if err := idx.store.Ping(); err != nil {
		health["database"] = "down"
		GinResponse(c, http.StatusInternalServerError, health)
		return
	}

	GinResponse(c, http.StatusOK, health)
}

func (idx *Indexer) handleIndex(c *gin.Context) {
	authToken, found := c.GetQuery("auth_token")
	if !found || len(authToken) < 1 {
		GinResponse(c, http.StatusUnauthorized, fmt.Errorf("auth_token is required"))
		return
	}

	if authToken != idx.conf.Indexer.Token {
		GinResponse(c, http.StatusUnauthorized, fmt.Errorf("auth_token is required"))
	}

	scanRange, found := c.GetQuery("scan")
	if !found || len(scanRange) < 1 {
		GinResponse(c, http.StatusBadRequest, fmt.Errorf("scan range is required"))
		return
	}

	start, end, err := utils.ParseScanQuery(scanRange, idx.latestBlock)
	if err != nil {
		GinResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// add block numbers to jobs queue channel
	go func(s, e int64) {
		for i := s; i <= e; i++ {
			idx.jobs <- i
		}
	}(start, end)

	GinResponse(c, http.StatusOK, "command executed")
}

// Response is the struct we return the response
type Response struct {
	Status  int         `json:"status"`
	Payload interface{} `json:"payload"`
}

func GinResponse(c *gin.Context, status int, payload interface{}) {
	response := Response{
		Status:  status,
		Payload: payload,
	}

	c.Header("Content-Type", "application/json")
	c.Status(status)

	c.JSON(status, response)
}
