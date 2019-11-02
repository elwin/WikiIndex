package app

import "github.com/gin-gonic/gin"

func (app *App) Serve() error {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "all good",
		})
	})
	
	r.GET("/size", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"size": app.index.Size(),
		})
	})

	r.GET("/page/:title", func(c *gin.Context) {
		title := c.Param("title")
		p, err := app.index.Get(title)
		if err != nil {
			c.JSON(404, gin.H{
				"error": err.Error(),
			})
			return
		}

		references := make([]string, 0)
		for reference := range p.References {
			references = append(references, reference.Title)
		}

		c.JSON(200, gin.H{
			"title": p.Title,
			"references": references,
		})
	})

	return r.Run(":8080")
}

