// package server contains a server for Hussar, which can be started with the function
// func Server()
package server

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/TblKVANYA/hussarGame/server/datatypes"
	"github.com/TblKVANYA/hussarGame/server/handler"
	"github.com/TblKVANYA/hussarGame/server/processor"
)

// func Server starts a server for Hussar
func Server() {
	str, err := getIP()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(str)
	str += ":8088"

	listener, err := net.Listen("tcp", str)
	if err != nil {
		log.Fatal(err)
	}

	var done [3]chan struct{}
	for i := 0; i < 3; i++ {
		done[i] = make(chan struct{})
	}

	var tunnels [3]datatypes.Tunnel
	for i := 0; i < 3; i++ {
		tunnels[i] = datatypes.TunnelInit()
	}

	go processor.Processor(tunnels)

	for i := 0; i < 3; i++ {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handler.HandleConn(conn, datatypes.Player(i), tunnels[i], done[i])
	}
	for i := 0; i < 3; i++ {
		<-done[i]
	}
}

// getIP returns IP in local net
func getIP() (string, error) {
	interSlice, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, inter := range interSlice {
		addrSlice, err := inter.Addrs()
		if err != nil {
			return "", nil
		}
		for _, addr := range addrSlice {
			str := addr.String()
			if str[:3] == "192" {
				return strings.Split(str, "/")[0], nil
			}
		}
	}
	return "", errors.New("ip: uppropriate IP is not found. Check your wi-fi connection")
}
