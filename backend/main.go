package main

import (
	"log"
	"os"
	"path/filepath"

	"gocroot/config"

	"github.com/gofiber/fiber/v2/middleware/cors"

	"gocroot/url"

	"github.com/gofiber/fiber/v2"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	config.InitDB()

	publicPath := filepath.Join(cwd, "frontend", "public")
	if _, err := os.Stat(publicPath); os.IsNotExist(err) {
		publicPath = filepath.Join(cwd, "..", "frontend", "public")
	}

	site := fiber.New(config.Iteung)
	site.Use(cors.New(config.Cors))
	site.Static("/", publicPath)
	url.Web(site)
	log.Fatal(site.Listen(config.IPPort))
}
