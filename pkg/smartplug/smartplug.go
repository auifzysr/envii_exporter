package smartplug

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
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

type client struct {
	addr    string
	timeout time.Duration
}

func encrypt(s string) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, struct{ uint32 }{uint32(len(s))})
	if err != nil {
		log.Println("failed: ", err)
		return nil
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

func New(address string) *client {
	return &client{
		addr:    address,
		timeout: time.Second * time.Duration(10),
	}
}

func (c *client) Dump() {
	fmt.Println("address: ", c.addr)
	fmt.Println("timeout: ", c.timeout)
}

func (c *client) Send(command string) {
	//d := net.Dialer{Timeout: c.timeout}
	//	conn, err := d.Dial("tcp", c.addr)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", c.addr)
	if err != nil {
		log.Println("failed to resolve tcp addr: ", err)
		return
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Println("failed to create conn: ", err)
		return
	}
	defer conn.Close()

	_, err = conn.Write(encrypt(command))
	if err != nil {
		fmt.Println("failed to write conn: ", err)
		return
	}
	result := [2048]byte{}
	_, err = conn.Read(result[0:])

	if err != nil {
		fmt.Println("failed to read conn: ", err)
		return
	}
	fmt.Println(decrypt(result[0:]))
}
