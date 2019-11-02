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
			"size": app.Count,
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

	r.GET("/path", func(c *gin.Context) {

		from, err := app.index.Get(c.Query("from"))
		if err != nil {
			c.JSON(404, gin.H{
				"error": err.Error(),
			})
			return
		}

		to, err := app.index.Get(c.Query("to"))
		if err != nil {
			c.JSON(404, gin.H{
				"error": err.Error(),
			})
			return
		}

		cost, err := app.index.Path(from, to)
		if err != nil {
			c.JSON(404, gin.H{
				"error": err.Error(),
			})
			return
		}


		c.JSON(200, gin.H{
			"from": from.Title,
			"to": to.Title,
			"cost": cost,
		})
	})

	r.GET("/longest", func(c *gin.Context) {

		from, err := app.index.Get(c.Query("from"))
		if err != nil {
			c.JSON(404, gin.H{
				"error": err.Error(),
			})
			return
		}


		to, cost := app.index.LongestPath(from)
		if to == nil {
			c.JSON(404, gin.H{
				"error": "no path found :(",
			})
			return
		}


		c.JSON(200, gin.H{
			"from": from.Title,
			"to":   to.Title,
			"cost": cost,
		})
	})

	r.GET("/loooongest", func(c *gin.Context) {

		from, to, cost := app.index.LongestTotalPath()

		c.JSON(200, gin.H{
			"from": from.Title,
			"to":   to.Title,
			"cost": cost,
		})
	})

	return r.Run(":8080")
}

