package main

import (
	"context"
	"log"
	"time"

	"github.com/ildarusmanov/go-up/app"
)

func HelloFactory(ctx context.Context) (app.Service, error) {
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
}
