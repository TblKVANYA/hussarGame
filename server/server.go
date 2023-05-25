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

	// start the listener
	listener, err := net.Listen("tcp", str)
	if err != nil {
		log.Fatal(err)
	}

	// init chans
	var done []chan struct{}
	done = append(done, make(chan struct{}))

	var tunnels []datatypes.Tunnel
	tunnels = append(tunnels, datatypes.TunnelInit())

	// some dances to get number of total players from someone who was the first to join
	numberChan := make(chan int32)
	conn, err := listener.Accept()
	if err != nil {
		log.Fatal(err)
	}
	go handler.HandleFirstConn(conn, tunnels[0], done[0], numberChan)

	N := <-numberChan

	// append some chans
	for i := int32(1); i < N; i++ {
		done = append(done, make(chan struct{}))
		tunnels = append(tunnels, datatypes.TunnelInit())
	}

	// start "processor"
	go processor.Processor(tunnels, N)

	// start N-1 connections
	for i := int32(1); i < N; i++ {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handler.HandleConn(conn, datatypes.Player(i), tunnels[i], done[i], N)
	}

	// wait for the end.
	for i := int32(0); i < N; i++ {
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
