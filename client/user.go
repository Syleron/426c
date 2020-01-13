package main

import "net"

type User struct {
	Username  string
	Listener  net.Conn
	Connected bool
}
