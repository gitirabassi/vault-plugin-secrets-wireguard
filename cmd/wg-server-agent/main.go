package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/hashicorp/vault/api"
	echo "github.com/labstack/echo/v4"
	"github.com/mitchellh/mapstructure"
	wgquick "github.com/nmiculinic/wg-quick-go"
	"github.com/sirupsen/logrus"
)

const (
	defualtInterface = "wg0"
)

var (
	serverCredsPath string
	webookSecret    string
	logger          = logrus.New()
)

// Cred is used for configuring a wireguard server
type Cred struct {
	Conf          string `mapstructure:"conf"`
	WebhookSecret string `mapstructure:"webhook_secret"`
}

func init() {
	flag.StringVar(&serverCredsPath, "creds-path", "wireguard/server-creds/default", "This is the path at which to find the credentials for the server that is going to be configured")
}

func main() {
	flag.Parse()

	initialized := false
	for !initialized {
		creds, err := getServerCred()
		if err != nil {
			log.Println("[init] can't get webhook token: ", err)
			time.Sleep(10 * time.Second)
			continue
		}
		webookSecret = creds.WebhookSecret
		err = wgQuickUp(creds.Conf)
		if err != nil {
			log.Println("[init] can't apply current configuration: ", err)
			time.Sleep(10 * time.Second)
			continue
		}
		initialized = true
	}
	e := echo.New()
	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})
	e.POST("/webhook", handleWebook)
	e.Logger.Fatal(e.Start(":51821"))
}

func getServerCred() (*Cred, error) {
	cli, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	secret, err := cli.Logical().Read(serverCredsPath)
	if err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, fmt.Errorf("secret is nil")
	}
	if secret.Data == nil {
		return nil, fmt.Errorf("secret.Data is nil")
	}
	serverCred := &Cred{}
	err = mapstructure.Decode(secret.Data, serverCred)
	if err != nil {
		return nil, err
	}

	return serverCred, nil
}

type WebhookAuth struct {
	Token string `json:"token" form:"token"`
}

func handleWebook(c echo.Context) error {
	auth := new(WebhookAuth)
	if err := c.Bind(auth); err != nil {
		return c.String(http.StatusBadRequest, "Couldn't parse tokes neither as json or form payload")
	}
	if auth.Token != webookSecret {
		return c.String(http.StatusUnauthorized, "Provided token doesn't match")
	}
	creds, err := getServerCred()
	if err != nil {
		log.Println("[init] can't get webhook token: ", err)
		c.String(http.StatusOK, "OK")
	}
	err = wgQuickSync(creds.Conf)
	if err != nil {
		log.Println("[init] can't apply current configuration: ", err)
		c.String(http.StatusOK, "OK")
	}
	return c.String(http.StatusOK, "OK")
}

func wgQuickUp(stringConf string) error {
	conf := &wgquick.Config{}
	err := conf.UnmarshalText([]byte(stringConf))
	if err != nil {
		return err
	}
	return wgquick.Up(conf, defualtInterface, logger)
}

func wgQuickSync(stringConf string) error {
	conf := &wgquick.Config{}
	err := conf.UnmarshalText([]byte(stringConf))
	if err != nil {
		return err
	}
	return wgquick.Sync(conf, defualtInterface, logger)
}
