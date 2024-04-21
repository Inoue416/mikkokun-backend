package handler

// import (
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// )

// type MikkokuPostType struct {
// 	TargetName string `json:"targetName"`
// 	UUID       string `json:"uuid"`
// }

// func mikkoku(targetName string, c *gin.Context) bool {
// 	var json MikkokuPostType
// 	if err := c.ShouldBind(&json); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"message": "そのユーザーは存在しません。"})
// 		return false
// 	}
// 	c.JSON(http.StatusOK, gin.H{"message": "密告しました。"})
// 	return true
// }
