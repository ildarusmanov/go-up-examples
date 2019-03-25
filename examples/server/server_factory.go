package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ildarusmanov/go-up/app"
)

func ServerFactory(ctx context.Context) (app.Service, error) {
	name, _ := ctx.Value("Application").(*app.Application).GetConfig("appName")

	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, fmt.Sprintf("Welcome to %s", name))
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	return srv, nil
}
