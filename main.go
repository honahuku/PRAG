package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/:host/auth/sign_in", func(c *gin.Context) {
		host := c.Param("host")

		uuid, err := generateUUID()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate UUID"})
			return
		}

		targetURL := fmt.Sprintf("https://%s/miauth/%s?name=PRAG&permission=read:account,write:notes", host, uuid)

		resp, err := http.Get(targetURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data"})
			return
		}
		defer resp.Body.Close()

		c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
	})

	r.Run()
}

func generateUUID() (string, error) {
	buffer := make([]byte, 16)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(buffer), nil
}
