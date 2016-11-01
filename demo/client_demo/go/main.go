// demo2 project main.go
package main

import (
	"flag"
	"fmt"
	"os"
	//"os/signal"

	zmq "github.com/alecthomas/gozmq"
)

func main() {
	var server_addr string
	var event_type string
	flag.StringVar(&server_addr, "zmq-server-addr", "tcp://127.0.0.1:8341", "zmq-server-addr for receving chaincode event")
	flag.StringVar(&event_type, "event_type", "", "zmq-server-addr for receving chaincode event")
	flag.Parse()

	fmt.Printf("Event server_addr: %v   event_type:%v \n", server_addr, event_type)

	context, _ := zmq.NewContext()
	socket, err := context.NewSocket(zmq.SUB)

	if err != nil {
		fmt.Println("1: ", err.Error())
		os.Exit(1)
	}

	err = socket.Connect(server_addr)
	socket.SetSubscribe(event_type)

	if err != nil {
		fmt.Println("2: ", err.Error())
		os.Exit(1)
	}

	for {
		println("before Event:")
		msg, _ := socket.Recv(0)
		println("Event:", string(msg))
	}

	println("try 2 close")
	socket.Disconnect(server_addr)
	socket.Close()
}
