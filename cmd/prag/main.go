package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"github.com/honahuku/prag/pkg/api/misskey"
)

func main() {
	certFile := "/etc/letsencrypt/live/misskey.sda1.net.prag.social/fullchain.pem"
	keyFile := "/etc/letsencrypt/live/misskey.sda1.net.prag.social/privkey.pem"

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		host := c.Request.Host
		hostParts := strings.Split(host, ".")
		if len(hostParts) < 3 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid host",
			})
			return
		}
		
		subdomainParts := hostParts[:len(hostParts)-2]
		subdomain := strings.Join(subdomainParts, ".")
		

		// 暫定的にすべてのリクエストがMastdonクライアントからMisskeyサーバーへのリクエストだとする
		// TODO: サーバーがどのActivityPub実装なのかを判定する必要がある
		// 暫定的にりなっくすきーでテストをしている
		// TODO: インスタンス追加の処理が外部から出来るようにする
		// インスタンス情報とかを同一インスタンス内にredisを立ててそこから取る？
		if subdomain == "misskey.sda1.net" {
			misskeyClient := misskey.NewMisskeyClient(subdomain)
			r.GET("/auth/sign_in", misskeyClient.SignInHandler())
		
			c.JSON(http.StatusOK, gin.H{
				"message": fmt.Sprintf("This request would be forwarded to %s server", subdomain),
			})
		} else {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Invalid subdomain",
			})
		}	
	})

	r.RunTLS(":443", certFile, keyFile)
}
