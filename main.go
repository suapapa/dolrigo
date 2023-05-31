package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/idtoken"
)

type ClientChan chan string

var (
	subPath     string
	templateDir string
)

func main() {
	flag.StringVar(&subPath, "p", "/", "subpath")                 // dolrigo
	flag.StringVar(&templateDir, "t", "template", "template dir") // -t /template
	flag.Parse()

	if !strings.HasSuffix(subPath, "/") {
		subPath += "/"
	}

	r := gin.Default()

	r.LoadHTMLGlob(path.Join(templateDir, "*.html"))
	r.GET(subPath+"ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong from dolrigo",
		})
	})
	r.GET(subPath+":gid/join", func(c *gin.Context) {
		gid := c.Param("gid")
		type JoinData struct {
			GID      string
			ClientID string
			SubPath  string
		}
		c.HTML(http.StatusOK, "join.html", &JoinData{
			GID:      gid,
			SubPath:  subPath,
			ClientID: os.Getenv("CLIENT_ID"),
		})
	})
	r.POST(subPath+":gid/login", func(c *gin.Context) {
		p, err := idtoken.Validate(c.Request.Context(), c.PostForm("credential"), "")
		if err != nil {
			c.String(500, "Internal Server Error 2")
			return
		}

		for k, v := range p.Claims {
			log.Println(k, v)
		}
		user := &Candidate{
			Name:  p.Claims["name"].(string),
			EMail: p.Claims["email"].(string),
			Photo: p.Claims["picture"].(string),
		}

		gid := c.Param("gid")
		if _, ok := Games[gid]; !ok {
			Games[gid] = NewGame()
		}
		game := Games[gid]
		game.AddCandidate(user)
	})
	r.GET(subPath+":gid/candidates", func(c *gin.Context) {
		gid := c.Param("gid")
		if _, ok := Games[gid]; !ok {
			Games[gid] = NewGame()
		}
		game := Games[gid]
		c.JSON(200, game)
	})
	r.StaticFS(subPath+":gid/", http.Dir("/static"))

	r.Run(":8080") // listen and serve on
}
