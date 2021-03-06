package zmq_wrapper

import (
	"github.com/op/go-logging"

	zmq "github.com/alecthomas/gozmq"
)

var wrapperLogger = logging.MustGetLogger("zmq_wrapper")

func InitZMQ(server_addr string) (*zmq.Socket, *string) {
	context, err := zmq.NewContext()
	if err != nil {
		err_str := "ZMQ can't create context!"
		return nil, &err_str
	}
	socket, err := context.NewSocket(zmq.PUB)
	if err != nil {
		err_str := "ZMQ can't create socket!"
		return nil, &err_str
	}
	err = socket.Bind(server_addr)
	if err != nil {
		err = socket.Close()
		if err != nil {
			err_str := "ZMQ can't close socket!"
			return nil, &err_str
		}
		err_str := "ZMQ can't bind the the address:" + server_addr
		return nil, &err_str
	}
	return socket, nil
}

func SendMsg(s *zmq.Socket, header string, content []byte) {
	wrapperLogger.Debugf("[%v]:%v", header, string(content))
	if s != nil {
		s.Send([]byte(header), zmq.SNDMORE)
		s.Send(content, 0)
	}
}
