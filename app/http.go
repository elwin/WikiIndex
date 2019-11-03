package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func (a *App) Serve() error {
	r := gin.Default()

	r.GET("/", a.Root())
	r.GET("/page/:title", a.Page())
	r.GET("/path", a.Path())
	r.GET("/longest", a.Longest())
	r.GET("/loooongest", a.LongestOverall())

	return r.Run(":8080")
}

func (a *App) Root() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "all good",
			"size":    a.Count,
		})
	}
}

func (a *App) Page() gin.HandlerFunc {
	return func(c *gin.Context) {
		title := c.Param("title")
		p, ok := a.Index.Get(title)
		if !ok {
			c.JSON(404, gin.H{
				"error": fmt.Sprintf("page '%s' not found", title),
			})
			return
		}

		referencesTo := make([]string, 0)
		for _, reference := range p.ReferencesTo() {
			referencesTo = append(referencesTo, reference.Title())
		}

		referencedBy := make([]string, 0)
		for _, reference := range p.ReferencedBy() {
			referencedBy = append(referencedBy, reference.Title())
		}

		c.JSON(200, gin.H{
			"title":        p.Title(),
			"referencesTo": referencesTo,
			"referencedBy": referencedBy,
		})
	}
}

func (a *App) Path() gin.HandlerFunc {
	return func(c *gin.Context) {
		from, ok := a.Index.Get(c.Query("from"))
		if !ok {
			c.JSON(404, gin.H{
				"error": fmt.Sprintf("page '%s' not found", from),
			})
			return
		}

		to, ok := a.Index.Get(c.Query("to"))
		if !ok {
			c.JSON(404, gin.H{
				"error": fmt.Sprintf("page '%s' not found", to),
			})
			return
		}

		cost, err := a.Index.Path(from, to)
		if err != nil {
			c.JSON(404, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"from": from.Title(),
			"to":   to.Title(),
			"cost": cost,
		})
	}
}

func (a *App) Longest() gin.HandlerFunc {
	return func(c *gin.Context) {
		from, ok := a.Index.Get(c.Query("from"))
		to, ok := a.Index.Get(c.Query("to"))
		if !ok {
			c.JSON(404, gin.H{
				"error": fmt.Sprintf("page '%s' not found", from),
			})
			return
		}

		to, cost := a.Index.LongestPath(from)
		if to == nil {
			c.JSON(404, gin.H{
				"error": "no path found :(",
			})
			return
		}

		c.JSON(200, gin.H{
			"from": from.Title(),
			"to":   to.Title(),
			"cost": cost,
		})
	}
}

func (a *App) LongestOverall() gin.HandlerFunc {
	return func(c *gin.Context) {
		from, to, cost := a.Index.LongestTotalPath()

		c.JSON(200, gin.H{
			"from": from.Title(),
			"to":   to.Title(),
			"cost": cost,
		})
	}
}
