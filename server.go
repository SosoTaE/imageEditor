package main

import (
	editor "media/imageEditor"
	"net/http"
	"github.com/gin-gonic/gin"
	"media/types"
)


func main() {
	r := gin.Default()

	r.POST("/api/v1/editor", func(ctx *gin.Context) {
		var jsonData types.JsonExample

		if err := ctx.BindJSON(&jsonData); err != nil {
			// JSON did not match expected format or required fields are missing
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		result, _ := editor.Edit(jsonData)

		
		ctx.JSON(http.StatusOK, gin.H{"array": result})
	})

	r.Run("0.0.0.0:9000")
}
