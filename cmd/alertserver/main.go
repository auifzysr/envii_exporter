package main

import (
	"log"

	"github.com/auifzysr/envii_exporter/pkg/alertserver"
	"github.com/auifzysr/envii_exporter/pkg/smartplug"
)

func main() {
	p := smartplug.New("192.168.11.43:9999")

	a := alertserver.NewAlertMux()
	a.AlertHandle(`{"state":"alerting","tags":{"threshold":"on"}}`, func() {
		p.Send(smartplug.ON)
	})
	a.AlertHandle(`{"state":"alerting","tags":{"threshold":"off"}}`, func() {
		p.Send(smartplug.OFF)
	})

	log.Fatal(alertserver.ListenAndServe(":32000", a))
}
