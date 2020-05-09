package main

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func main() {
	webhookSecret := ""
	for webhookSecret == "" {
		webhookSecret, err := getWebhookSecret()
		if err != nil {
			log.Prinln("[init] can't get webhook token: ", err)
		}
		time.Sleep(10 * time.Second)
	}
	e := echo.New()
	e.GET("/webhook", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!"+webhookSecret)
	})
	e.Logger.Fatal(e.Start(":1323"))
}

func getWebhookSecret() (string, error) {
	return "", nil
}
