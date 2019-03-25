package main

import (
	"context"

	"github.com/ildarusmanov/go-up/app"
)

func PrinterFactory(ctx context.Context) (app.Service, error) {
	return &Printer{}, nil
}
