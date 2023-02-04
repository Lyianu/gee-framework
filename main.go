package main

import (
	"gee"
	"log"
	"net/http"
	"time"
)

func main() {
	r := gee.New()
	r.Use(logger)
	v1 := r.Group("/v1")
	{
		v1.GET("/", func(c *gee.Context) {
			c.HTML(http.StatusOK, "<h1>Welcome</h1>")
		})

		v1.GET("/hello", func(c *gee.Context) {
			c.String(http.StatusOK, "hello!")
		})

		api := v1.Group("/api")
		{
			api.GET("/", func(c *gee.Context) {
				c.String(http.StatusOK, "Welcome to api!")
			})

			api.GET("/", func(c *gee.Context) {
				c.String(http.StatusOK, c.Query("s"))
			})
		}
		//api.Use(apiLogger)
	}
	r.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Welcome!</h1><p>S="+c.Query("s")+"</p>")
	})

	r.GET("/hello", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{"1": 1, "s": c.Query("s")})
	})

	r.GET("/hello/this", nil)
	r.POST("/hello", nil)

	r.Router().Traverse()
	r.Run(":9999")
}

func logger(c *gee.Context) {
	t := time.Now()
	c.Next()
	timeInteval := time.Since(t)

	log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, timeInteval)
}
