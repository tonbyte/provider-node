package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func registerHandlers(e *echo.Echo, h *handler) {
	postCors := middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.POST},
	})
	getCors := middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET},
	})

	dapp := e.Group("/v1/provider")
	dapp.Use(middleware.CORS())
	dapp.GET("/status", h.Status, getCors)
	dapp.POST("/uploadFile", h.UploadFile, postCors)
}
