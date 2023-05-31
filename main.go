package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/idtoken"
)

var (
	survyeID = "1P_Myt5-7aRncmZ2AJwC0UwhkmEAEqxCdJpUxhjmp_SM"
)

func main() {

	r := gin.Default()
	r.LoadHTMLGlob("template/*.html")
	r.GET("/dolrigo/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong from dolrigo",
		})
	})
	// r.GET("/", func(c *gin.Context) {
	// 	c.String(404, "Not Found")
	// })
	r.GET("/dolrigo", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	r.POST("/dolrigo/login", func(c *gin.Context) {
		p, err := idtoken.Validate(c.Request.Context(), c.PostForm("credential"), "")
		if err != nil {
			c.String(500, "Internal Server Error 2")
			return
		}

		for k, v := range p.Claims {
			log.Println(k, v)
		}
	})

	r.Run(":8080") // listen and serve on

	// // Get the Google Survey results.
	// resp, err := http.Get(fmt.Sprintf("https://forms.googleapis.com/v1/forms/%s/responses", survyeID))
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// defer resp.Body.Close()
	// io.Copy(os.Stdout, resp.Body)

	// Decode the JSON response.
	// results := make([]map[string]interface{}, 0)
	// err = json.NewDecoder(resp.Body).Decode(&results)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// // Do something with the results.
	// for _, result := range results {
	// 	fmt.Println(result)
	// }
}
