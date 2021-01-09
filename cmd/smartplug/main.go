package main

import (
	"flag"
	"time"

	"github.com/auifzysr/envii_exporter/pkg/smartplug"
)

func main() {
	hostname := flag.String("hostname", "127.0.0.1", "hostname")
	port := flag.String("port", "9999", "port")
	timeout := flag.Int("timeout", 10, "timeout")
	flag.Parse()

	client := smartplug.New(&smartplug.Config{
		Hostname: *hostname,
		Port:     *port,
		Timeout:  time.Second * time.Duration(*timeout),
	})

	client.Dump()

}
