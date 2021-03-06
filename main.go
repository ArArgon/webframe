package main

import (
	"log"

	"github.com/ArArgon/webframe/lib"
)

func main() {
	engine := lib.New()

	engine.GET("/", func(c *lib.Context) {
		c.SetHeader("Location", "/info")
		c.SetStatusCode(301)
	})

	engine.GET("/info", func(c *lib.Context) {
		c.String(400, "no info")
	})

	engine.GET("/json", func(c *lib.Context) {
		c.JSON(400, lib.JSONObject{
			"Hello":  "world",
			"isOK":   true,
			"anList": []interface{}{1, "fine"},
		})
	})

	engine.GET("/json/:obj", func(c *lib.Context) {
		c.String(400, "Part match success: %v", c.Params)
	})

	engine.GET("/path2/*match", func(c *lib.Context) {
		c.String(400, "Catch-all match success: %v", c.Params)
	})

	engine.GET("/path3/:hybrid/*match", func(c *lib.Context) {
		c.String(400, "Hybrid matches success: %v", c.Params)
	})

	v1 := engine.CreateSubGroup("/v1")
	v1.GET("/path4/:hybrid/*match", func(c *lib.Context) {
		c.String(400, "Hybrid matches in control group success: %v", c.Params)
	})

	v1.GET("/panic", func(c *lib.Context) {
		panic("A deliberate panic trap")
	})

	v1.Static("/assets", "./static", false)

	log.Printf("[Main] Engine launching...")
	log.Fatal(engine.Run(":9000"))
}
