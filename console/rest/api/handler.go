package api

import (
	"encoding/json"
	"energi-challenge/console/rest/utils"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// handleHealthCheck checks database ping
func (a *API) handleHealthCheck(c *gin.Context) {
	health := map[string]interface{}{
		"store": "up",
	}

	if err := a.store.Ping(); err != nil {
		health["database"] = "down"
		c.JSON(http.StatusInternalServerError, health)
		return
	}

	c.JSON(http.StatusOK, health)
}

// handleIndexerCommand to connect to the indexer via the rest APIs
func (a *API) handleIndexerCommand(c *gin.Context) {
	scanRange, found := c.GetQuery("scan")
	if !found || len(scanRange) < 1 {
		a.response(c, http.StatusBadRequest, "scan parameter is required")
		return
	}

	resp, err := http.Get(fmt.Sprintf("%s%s?auth_token=%s&scan=%s", a.conf.Indexer.Host, a.conf.Indexer.Address, a.conf.Indexer.Token, scanRange))
	if err != nil {
		a.response(c, http.StatusBadGateway, err.Error())
		return
	}

	if resp.StatusCode != http.StatusOK {
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(resp.Body)

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			a.response(c, http.StatusBadRequest, err.Error())
			return
		}
		var body map[string]any
		if err := json.Unmarshal(b, &body); err != nil {
			a.response(c, http.StatusBadRequest, err.Error())
			return
		}

		a.response(c, http.StatusBadRequest, body)
		return
	}

	a.response(c, http.StatusOK, "indexer command executed")
}

// handleGetLatestBlock returns the latest block and all associated transactions
func (a *API) handleGetLatestBlock(c *gin.Context) {

	block, err := a.repository.GetLatestBlock(c.Request.Context())
	if err != nil {
		log.Printf("failed to get the latest block: %s", err.Error())
		a.response(c, http.StatusInternalServerError, err.Error())
		return
	}

	a.response(c, http.StatusOK, block)
}

// handleGetBlock returns the provided block number information
func (a *API) handleGetBlock(c *gin.Context) {
	id := c.Param("id")
	num, err := strconv.ParseInt(id, 10, 0)
	if err != nil {
		a.response(c, http.StatusBadRequest, err.Error())
		return
	}

	// is negative
	if num < 0 {
		num = num * -1
	}

	block, err := a.repository.GetBlock(c.Request.Context(), num)
	if err != nil {
		a.response(c, http.StatusInternalServerError, err.Error())
		return
	}

	a.response(c, http.StatusOK, block)
}

// handleGetStats returns the latest block and the transaction stats
func (a *API) handleGetStats(c *gin.Context) {
	ctx := c.Request.Context()

	var latest int64 = 0
	block, err := a.repository.GetLatestBlock(ctx)
	if err == nil {
		latest = block.Number
	}

	stats, err := a.repository.GetStats(ctx, 0, latest)
	if err != nil {
		a.response(c, http.StatusInternalServerError, err.Error())
		return
	}

	a.response(c, http.StatusOK, stats)
}

// handleGetRangeStats returns the stats for the provided block range
func (a *API) handleGetRangeStats(c *gin.Context) {
	ctx := c.Request.Context()

	interval := c.Param("range")
	start, end, err := utils.ParseRange(interval)
	if err != nil {
		a.response(c, http.StatusBadRequest, err.Error())
		return
	}

	stats, err := a.repository.GetStats(ctx, start, end)
	if err != nil {
		a.response(c, http.StatusInternalServerError, err.Error())
		return
	}

	a.response(c, http.StatusOK, stats)
}

// handleGetLatestTx return the latest transaction information
func (a *API) handleGetLatestTx(c *gin.Context) {
	tx, err := a.repository.GetLatestTx(c.Request.Context())
	if err != nil {
		a.response(c, http.StatusInternalServerError, err.Error())
		return
	}
	a.response(c, http.StatusOK, tx)
}

// handleGetTx returns the provided transaction hash information
func (a *API) handleGetTx(c *gin.Context) {
	hash := c.Param("hash")
	tx, err := a.repository.GetTx(c.Request.Context(), hash)
	if err != nil {
		a.response(c, http.StatusInternalServerError, err.Error())
		return
	}
	a.response(c, http.StatusOK, tx)
}
