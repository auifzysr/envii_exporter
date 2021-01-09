package smartplug

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

const (
	INFO      = `{"system":{"get_sysinfo":{}}}`
	ON        = `{"system":{"set_relay_state":{"state":1}}}`
	OFF       = `{"system":{"set_relay_state":{"state":0}}}`
	LEDOFF    = `{"system":{"set_led_off":{"off":1}}}`
	LEDON     = `{"system":{"set_led_off":{"off":0}}}`
	CLOUDINFO = `{"cnCloud":{"get_info":{}}}`
	WLANSCAN  = `{"netif":{"get_scaninfo":{"refresh":0}}}`
	TIME      = `{"time":{"get_time":{}}}`
	SCHEDULE  = `{"schedule":{"get_rules":{}}}`
	COUNTDOWN = `{"count_down":{"get_rules":{}}}`
	ANTITHEFT = `{"anti_theft":{"get_rules":{}}}`
	REBOOT    = `{"system":{"reboot":{"delay":1}}}`
	RESET     = `{"system":{"reset":{"delay":1}}}`
	ENERGY    = `{"emeter":{"get_realtime":{}}}`
)

var key = 171

type Config struct {
	Hostname string
	Port     string
	Timeout  time.Duration
}

type client struct {
	hostname string
	port     string
	timeout  time.Duration
}

func encrypt(s string) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, struct{ uint32 }{uint32(len(s))})
	if err != nil {
		fmt.Println("failed: ", err)
	}
	result := []byte(buf.Bytes())
	for _, c := range []byte(s) {
		a := key ^ int(c)
		key = a
		result = append(result, byte(a))
	}
	return result
}

func decrypt(s []byte) string {
	result := ``
	for _, c := range s {
		a := key ^ int(c)
		key = int(c)
		result += string(a)
	}
	return result
}

func New(cfg *Config) *client {
	return &client{
		hostname: cfg.Hostname,
		port:     cfg.Port,
		timeout:  cfg.Timeout,
	}
}

func (c *client) Dump() {
	fmt.Println("hostname: ", c.hostname)
	fmt.Println("port: ", c.port)
	fmt.Println("timeout: ", c.timeout)
}

func (c *client) Send(command string) {
	d := net.Dialer{Timeout: time.Duration(c.timeout) * time.Second}
	conn, err := d.Dial("tcp", fmt.Sprintf("%s:%s", c.hostname, c.port))
	if err != nil {
		fmt.Println("failed to create conn: ", err)
	}
	defer conn.Close()

	_, err = conn.Write(encrypt(command))
	if err != nil {
		fmt.Println("failed to write conn: ", err)
	}
	result := [2048]byte{}
	_, err = conn.Read(result[0:])

	if err != nil {
		fmt.Println("failed to read conn: ", err)
	}
	fmt.Println(decrypt(result[0:]))
}
