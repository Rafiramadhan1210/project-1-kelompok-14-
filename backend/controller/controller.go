package controller

import (
	"gocroot/helper"

	"github.com/gofiber/fiber/v2"
)

func Homepage(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON("Welcome to Go Croot!!!")
}

func GetIPServer(c *fiber.Ctx) error {
	ipaddr := helper.GetIPaddress()
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"ip_address": ipaddr})
}
