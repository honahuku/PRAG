package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"github.com/google/uuid"
)

type MisskeyClient struct {
	Host string
}

func NewMisskeyClient(host string) *MisskeyClient {
	return &MisskeyClient{Host: host}
}

func (c *MisskeyClient) SignInHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// UUIDを生成
		uuid := uuid.New().String()

		// ログインURLを生成
		loginURL := fmt.Sprintf("https://%s/%s/miauth/?name=PRAG&permission=read:account,write:notes", c.Host, uuid)

		// HTTPクライアントを作成
		client := &http.Client{}

		// リクエストを作成
		req, err := http.NewRequest("GET", loginURL, nil)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to create request",
			})
			return
		}

		// リクエストを送信
		resp, err := client.Do(req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to send request",
			})
			return
		}
		defer resp.Body.Close()

		// レスポンスボディを読み込む
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to read response",
			})
			return
		}

		// レスポンスボディをクライアントに返す
		ctx.String(http.StatusOK, string(body))
	}
}
