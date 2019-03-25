package main

import "log"

type Printer struct{}

func (p *Printer) SayHello(name string) {
	log.Printf("Hello %s!\n", name)
}
