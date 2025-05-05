package api

import (
	"github.com/darcygail/ether-explorer/internal/store"
	"github.com/gin-gonic/gin"
)

func RunApi() {
	r := gin.Default()
	r.GET("/asset/:address", func(c *gin.Context) {
		addr := c.Param("address")
		ctx := c.Request.Context()
		asset, err := store.GetAsset(ctx, addr)
		if err != nil {
			c.JSON(404, gin.H{"error": "not found"})
			return
		}
		c.JSON(200, asset)
	})

	r.Run(":8080")
}