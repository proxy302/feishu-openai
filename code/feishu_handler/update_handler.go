package feishu_handler

import (
	"start-feishubot/aihandlers"

	"github.com/gin-gonic/gin"
)

type UpdateRequest struct {
	ID int `json:"id"`
}

func UpdateHandler(c *gin.Context) {
	var (
		updateRequest UpdateRequest
	)
	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	aihandlers.UpdateHandler(updateRequest.ID)
	c.JSON(200, gin.H{"message": "success"})
}
