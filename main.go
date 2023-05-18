package main

import (
	"crypto/ed25519"
	"strconv"

	"github.com/tonbyte/provider-node/config"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

var Version = "1"

func main() {
	e := echo.New()
	config.LoadConfig()

	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		log.Fatalf("generate keys error: %v", err)
		return
	}

	h := newHandler(pub, priv)

	registerHandlers(e, h)

	e.Logger.Fatal(e.Start(":" + strconv.Itoa(config.StorageConfig.Port)))
}
