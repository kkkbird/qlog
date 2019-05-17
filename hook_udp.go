package qlog

import (
	"net"
	"reflect"
)

const (
	keyUDPEnabled = "logger.udp.enabled"
	keyUDPLevel   = "logger.udp.level"
	keyUDPAddress = "logger.udp.address"
)

// TODO: define a udpLogger as writer, use concurrent module to write data
// type udpLogger struct {
// 	conn net.Conn
// }

// func newUDPLogger(addr string) (*udpLogger, error) {
// 	udpAddr, err := net.ResolveUDPAddr("udp", addr)

// 	if err != nil {
// 		return nil, err
// 	}

// 	conn, err := net.DialUDP("udp", nil, addr)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return &udpLogger{
// 		conn: conn,
// 	}, nil
// }

// UDPHook output message to udp
type UDPHook struct {
	BaseHook

	Address string
}

// Setup function for UDPHook
func (h *UDPHook) Setup() (err error) {
	h.baseSetup()

	h.Address = v.GetString(keyUDPAddress)

	udpAddr, err := net.ResolveUDPAddr("udp", h.Address)
	if err != nil {
		return err
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return err
	}

	h.writer = conn

	return nil
}

var _InitUDPHook = func() interface{} {
	cli.Bool(keyUDPEnabled, false, "logger.udp.enabled")
	cli.String(keyUDPLevel, "", "logger.udp.level") // DONOT set default level in pflag

	registerHook("udp", reflect.TypeOf(UDPHook{}))
	return nil
}()
