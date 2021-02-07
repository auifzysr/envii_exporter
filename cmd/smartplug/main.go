package main

import (
	"flag"
	"time"

	"github.com/auifzysr/envii_exporter/pkg/smartplug"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:9999", "hostname")
	flag.Parse()

	client := smartplug.New(*addr)

	client.Send(smartplug.INFO)
	time.Sleep(time.Second * 3)
	// BUG: causes "failed to read conn:  EOF"
	client.Send(smartplug.INFO)
	time.Sleep(time.Second * 3)
	client.Send(smartplug.INFO)
}
