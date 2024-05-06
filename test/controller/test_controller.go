// Create by code generator  2024-05:06 17:39:04.024
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Index(g *gin.Context) {
	data := gin.H{
		"code": 200,
	}
	g.JSON(http.StatusOK, data)
}
