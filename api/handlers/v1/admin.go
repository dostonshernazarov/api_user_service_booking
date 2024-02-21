package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *handlerV1) Admin(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{
		"message": "Welcome admin",
	})
}
