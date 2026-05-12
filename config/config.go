package config

import (
	"gocroot/helper"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

var _ = godotenv.Load()

var IPPort, Net = helper.GetAddress()

var Iteung = fiber.Config{
	Prefork:       false,
	CaseSensitive: true,
	StrictRouting: true,
	ServerHeader:  "GoCroot",
	AppName:       "Golang Change Root",
	Network:       Net,
}
