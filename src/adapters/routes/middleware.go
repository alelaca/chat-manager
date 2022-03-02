package routes

import (
	"github.com/gin-gonic/gin"
)

func configureCors(c *gin.Context) {
	// this is set here only for the static html frontend to work!!
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
}
