package main

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/html"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
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

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
			return
		}

		newBody := rewriteLinks(body, host)
		c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), newBody)
	})

	certFile := "/etc/letsencrypt/live/tacosync.io/fullchain.pem"
	keyFile := "/etc/letsencrypt/live/tacosync.io/privkey.pem"

	log.Fatal(r.RunTLS(":443", certFile, keyFile))
}

func generateUUID() (string, error) {
	buffer := make([]byte, 16)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(buffer), nil
}

func rewriteLinks(body []byte, host string) []byte {
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return body
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			// Consider 'a', 'link', and 'script' tags
			if n.Data == "a" || n.Data == "link" || n.Data == "script" {
				rewriteAttributes(n, host)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)

	var buf bytes.Buffer
	html.Render(&buf, doc)
	return buf.Bytes()
}

func rewriteAttributes(n *html.Node, host string) {
	attributes := []string{"href", "src"}
	for i, a := range n.Attr {
		for _, attr := range attributes {
			if a.Key == attr {
				// Rewrite absolute URLs
				if strings.HasPrefix(a.Val, "https://") || strings.HasPrefix(a.Val, "http://") {
					n.Attr[i].Val = strings.Replace(a.Val, "https://", fmt.Sprintf("https://%s/", host), 1)
					n.Attr[i].Val = strings.Replace(a.Val, "http://", fmt.Sprintf("http://%s/", host), 1)
				} else if !strings.HasPrefix(a.Val, "data:") && !strings.HasPrefix(a.Val, "javascript:") {
					// Prepend the host to relative URLs
					n.Attr[i].Val = fmt.Sprintf("https://%s%s", host, a.Val)
				}
			}
		}
	}
}

