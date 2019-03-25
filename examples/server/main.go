package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ildarusmanov/go-up/app"
)

func main() {
	log.Println("[+] Starting")

	sigchan := make(chan os.Signal, 1)

	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

	a := app.NewApplication(ctx, nil, nil)

	a.SetConfig("userName", "User")

	a.SetConfig("appName", "Test App")

	a.AddServiceFactory("printer", func(ctx context.Context) (app.Service, error) {
		return &Printer{}, nil
	})

	a.AddServiceFactory("hello", func(ctx context.Context) (app.Service, error) {
		go func() {
			for {
				select {
				case <-ctx.Done():
					log.Printf("[*] Hello exited with: %s", ctx.Err())
					return
				default:
					p, err := ctx.Value("Application").(*app.Application).GetService("printer")
					if err != nil {
						log.Println(err)
					} else {
						name, _ := ctx.Value("Application").(*app.Application).GetConfig("userName")

						p.(*Printer).SayHello(name)
					}
					time.Sleep(time.Second * 1)
				}
			}
		}()

		return nil, nil
	})

	a.AddServiceFactory("server", func(ctx context.Context) (app.Service, error) {
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
	})

	log.Println("[*] Started")

	<-sigchan

	stopServer(a)
	cancel()

	time.Sleep(time.Second * 5)

	log.Println("[x] Finished")

	os.Exit(0)
}

func stopServer(a *app.Application) {
	srv, err := a.GetService("server")

	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.(*http.Server).Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	log.Println("[*] Server stopped")
}

type Printer struct{}

func (p *Printer) SayHello(name string) {
	log.Printf("Hello %s!\n", name)
}
