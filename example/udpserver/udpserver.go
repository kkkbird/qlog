package main

import (
	"context"
	"fmt"
	"net"
)

const (
	address = ":31234"
)

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		fmt.Println(err)
		return
	}

	conn, err := net.ListenUDP("udp", udpAddr)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Listen at", address)

	ctx := context.Background()

	go func() {
		buf := make([]byte, 65535)
		for {
			len, _, err := conn.ReadFromUDP(buf)

			if err != nil {
				fmt.Println("ReadFromUDP err:", err)
				return
			}
			fmt.Println("len:", len, "buf:", string(buf[0:len]))
		}
	}()

	select {
	case <-ctx.Done():
		return
	}
}
