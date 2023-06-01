package main

import (
	"flag"
	"fmt"
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
	r.GET(subPath+"join/:gid", func(c *gin.Context) {
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
	r.POST(subPath+"login/:gid", func(c *gin.Context) {
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

		c.String(200, fmt.Sprintf("%s 돌림판에 참여하셨습니다.", user.Name))
	})
	r.DELETE(subPath+"candidates/:gid/:email", func(c *gin.Context) {
		gid := c.Param("gid")
		email := c.Param("email")
		if _, ok := Games[gid]; !ok {
			Games[gid] = NewGame()
		}
		game := Games[gid]
		game.RemoveCandidate(email)
		c.String(200, fmt.Sprintf("%s 돌림판에서 %s 이(가) 제외되었습니다.", gid, email))
	})
	r.GET(subPath+"candidates/:gid", func(c *gin.Context) {
		gid := c.Param("gid")
		if _, ok := Games[gid]; !ok {
			Games[gid] = NewGame()
		}
		game := Games[gid]
		c.JSON(200, game.Candidates)
	})
	r.StaticFS(subPath+"static", http.Dir("/static"))
	r.GET(subPath, func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "static/index.html")
	})

	r.Run(":8080") // listen and serve on
}
