package main

import (
	"net/http"

	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/sukso96100/covid19-push/fcm"
	// "io/ioutil"
)

type Token struct {
	Token string `json:"token" form:"token" query:"token"`
}

func Subscribe(c echo.Context) error {
	topic := c.QueryParam("topic")

	if topic != "stat" && topic != "news" {
		return c.String(http.StatusBadRequest, "That topic is not available")
	}

	token := new(Token)
	if err := c.Bind(token); err != nil {
		fmt.Println("%w", err)
		return c.String(http.StatusBadRequest, "Token form invalid!")
	}

	fcmApp := fcm.GetFCMApp()
	response, err := fcmApp.MsgClient.SubscribeToTopic(fcmApp.Ctx, []string{token.Token}, topic)
	if err != nil {
		fmt.Printf("%w", err)
		return c.String(http.StatusInternalServerError, "Subscribe Error")
	}
	return c.String(http.StatusOK, fmt.Sprintf("Subscribed!: %d", response.SuccessCount))

}

func Unubscribe(c echo.Context) error {
	topic := c.QueryParam("topic")

	if topic != "stat" && topic != "news" {
		return c.String(http.StatusBadRequest, "That topic is not available")
	}

	token := new(Token)
	if err := c.Bind(token); err != nil {
		fmt.Println("%w", err)
		return c.String(http.StatusBadRequest, "Token form invalid!")
	}

	fcmApp := fcm.GetFCMApp()
	response, err := fcmApp.MsgClient.UnsubscribeFromTopic(fcmApp.Ctx, []string{token.Token}, topic)
	if err != nil {
		fmt.Printf("%w", err)
		return c.String(http.StatusInternalServerError, "Unubscribe Error")
	}
	return c.String(http.StatusOK, fmt.Sprintf("Subscribed!: %d", response.SuccessCount))
}
