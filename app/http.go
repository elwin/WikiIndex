package app

import (
	"WikiIndex/database"
	"fmt"
	"github.com/flosch/pongo2"
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
		tpl := pongo2.Must(pongo2.FromFile("view/index.html"))
		tpl = tpl

		//c.JSON(200, gin.H{
		//	"message": "all good",
		//	"size":    a.Count,
		//})
		//
		//tpl.Execute(pongo2.Context{
		//
		//})

		err := tpl.ExecuteWriter(pongo2.Context{}, c.Writer)
		if err != nil {
			fmt.Println(err)
		}
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

		tpl := pongo2.Must(pongo2.FromFile("view/page.html"))
		err := tpl.ExecuteWriter(pongo2.Context{
			"title":        p.Title(),
			"slug":         p.Slug(),
			"referencesTo": p.ReferencesTo(),
			"referencedBy": p.ReferencedBy(),
		}, c.Writer)
		if err != nil {
			fmt.Println(err)
		}
		//
		//c.JSON(200, gin.H{
		//	"title":        p.Title(),
		//	"referencesTo": referencesTo,
		//	"referencedBy": referencedBy,
		//})
	}
}

func (a *App) Path() gin.HandlerFunc {
	//return func(c *gin.Context) {
	//
	//	tpl := pongo2.Must(pongo2.FromFile("view/path.html"))
	//	err := tpl.ExecuteWriter(pongo2.Context{}, c.Writer)
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//
	//}
	return func(c *gin.Context) {
		result := struct {
			Set   bool
			Error error
			Path  []database.Pageable
			Len   int
		}{}

		from, to := c.Query("from"), c.Query("to")
		if from != "" && to != "" {
			from, ok := a.Index.Get(from)
			if !ok {
				c.JSON(404, gin.H{
					"error": fmt.Sprintf("page '%s' not found", from),
				})
				return
			}

			to, ok := a.Index.Get(to)
			if !ok {
				c.JSON(404, gin.H{
					"error": fmt.Sprintf("page '%s' not found", to),
				})
				return
			}

			result.Path, result.Error = a.Index.Path(from, to)
			result.Len = len(result.Path)
			result.Set = true
		}

		tpl := pongo2.Must(pongo2.FromFile("view/path.html"))
		err := tpl.ExecuteWriter(pongo2.Context{"result": result}, c.Writer)
		if err != nil {
			fmt.Println(err)
		}

		//c.JSON(200, gin.H{
		//	"from": from.Title(),
		//	"to":   to.Title(),
		//	"cost": cost,
		//})
	}
}

func (a *App) Longest() gin.HandlerFunc {
	return func(c *gin.Context) {
		from, ok := a.Index.Get(c.Query("from"))
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
