// demo project main.go
package main

import (
	"fmt"
	"os"

	zmq "github.com/alecthomas/gozmq"
)

func main() {
	context, _ := zmq.NewContext()
	socket, err := context.NewSocket(zmq.PUB)

	if err != nil {
		fmt.Println("2: ", err.Error())
		os.Exit(1)
	}

	err = socket.Bind("tcp://127.0.0.1:8341")

	if err != nil {
		fmt.Println("2: ", err.Error())
		os.Exit(1)
	}
	for i := 0; i < 1000000; i++ {
		msg := fmt.Sprintf("msg %d", i)
		socket.Send([]byte(msg), 0)
		println("Sending", msg)
	}
}
