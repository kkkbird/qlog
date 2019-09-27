package qlog

import (
	"net"
	"reflect"

	"github.com/sirupsen/logrus"
)

const (
	keyUDPEnabled = "logger.udp.enabled"
	keyUDPLevel   = "logger.udp.level"
	keyUDPHost    = "logger.udp.host"
	keyUDPUUID    = "logger.udp.uuid"
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
	b BaseHook

	Name string // for new hook
	Host string
	UUID string
}

// Fire output message to hook writer
func (h *UDPHook) Fire(e *logrus.Entry) error {
	if len(h.UUID) > 0 {
		e.Data["uuid"] = h.UUID
		defer delete(e.Data, "uuid")
	}
	return h.b.Fire(e)
}

// Levels return all available debug level of a hook
func (h *UDPHook) Levels() []logrus.Level {
	return h.b.Levels()
}

// Setup function for UDPHook
func (h *UDPHook) Setup() (err error) {
	h.b.Name = h.Name // name is set by reflect
	h.b.baseSetup()

	h.UUID = v.GetString(keyUDPUUID)
	h.Host = v.GetString(keyUDPHost)

	udpAddr, err := net.ResolveUDPAddr("udp", h.Host)
	if err != nil {
		return err
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return err
	}

	h.b.writer = conn

	return nil
}

var _InitUDPHook = func() interface{} {
	cli.Bool(keyUDPEnabled, false, "logger.udp.enabled")
	cli.String(keyUDPLevel, "", "logger.udp.level") // DONOT set default level in pflag

	registerHook("udp", reflect.TypeOf(UDPHook{}))
	return nil
}()
