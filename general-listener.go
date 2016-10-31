/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package main

import (
	//"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/leochan007/fabric-listener/events/consumer"
	pb "github.com/leochan007/fabric-listener/protos"
	"github.com/op/go-logging"

	"github.com/leochan007/fabric-listener/zmq_wrapper"
)

type adapter struct {
	notfy              chan *pb.Event_Block
	rejected           chan *pb.Event_Rejection
	cEvent             chan *pb.Event_ChaincodeEvent
	listenToRejections bool
	chaincodeID        string
}

var listenerLogger = logging.MustGetLogger("general_listener")

//GetInterestedEvents implements consumer.EventAdapter interface for registering interested events
func (a *adapter) GetInterestedEvents() ([]*pb.Interest, error) {
	if a.chaincodeID != "" {
		return []*pb.Interest{
			{EventType: pb.EventType_BLOCK},
			{EventType: pb.EventType_REJECTION},
			{EventType: pb.EventType_CHAINCODE,
				RegInfo: &pb.Interest_ChaincodeRegInfo{
					ChaincodeRegInfo: &pb.ChaincodeReg{
						ChaincodeID: a.chaincodeID,
						EventName:   ""}}}}, nil
	}
	return []*pb.Interest{{EventType: pb.EventType_BLOCK}, {EventType: pb.EventType_REJECTION}}, nil
}

//Recv implements consumer.EventAdapter interface for receiving events
func (a *adapter) Recv(msg *pb.Event) (bool, error) {
	if o, e := msg.Event.(*pb.Event_Block); e {
		a.notfy <- o
		return true, nil
	}
	if o, e := msg.Event.(*pb.Event_Rejection); e {
		if a.listenToRejections {
			a.rejected <- o
		}
		return true, nil
	}
	if o, e := msg.Event.(*pb.Event_ChaincodeEvent); e {
		a.cEvent <- o
		return true, nil
	}
	return false, fmt.Errorf("Receive unkown type event: %v", msg)
}

//Disconnected implements consumer.EventAdapter interface for disconnecting
func (a *adapter) Disconnected(err error) {
	fmt.Printf("Disconnected...exiting\n")
	os.Exit(1)
}

func createEventClient(eventAddress string, listenToRejections bool, cid string) *adapter {
	var obcEHClient *consumer.EventsClient

	done := make(chan *pb.Event_Block)
	reject := make(chan *pb.Event_Rejection)
	adapter := &adapter{notfy: done, rejected: reject, listenToRejections: listenToRejections, chaincodeID: cid, cEvent: make(chan *pb.Event_ChaincodeEvent)}
	obcEHClient, _ = consumer.NewEventsClient(eventAddress, 5, adapter)
	if err := obcEHClient.Start(); err != nil {
		fmt.Printf("could not start chat %s\n", err)
		obcEHClient.Stop()
		return nil
	}

	return adapter
}

var break_out bool = false

func systemEventHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	s := <-c
	fmt.Println("Got signal:", s)
	break_out = true
}

const (
	notfy_filter    = "notfy"
	rejected_filter = "rejected"
	cEvent_filter   = "cEvent"
)

func main() {
	var eventAddress string
	var listenToRejections bool
	var chaincodeID string
	var server_addr string
	var sleetime int64
	flag.StringVar(&eventAddress, "events-address", "0.0.0.0:7053", "address of events server")
	flag.BoolVar(&listenToRejections, "listen-to-rejections", false, "whether to listen to rejection events")
	flag.StringVar(&chaincodeID, "events-from-chaincode", "", "listen to events from given chaincode")
	flag.StringVar(&server_addr, "zmq-server-addr", "tcp://0.0.0.0:8341", "zmq-server-addr for sending chaincode event")
	flag.Int64Var(&sleetime, "sleeptime", 1, "sleep time(Millisecond) for socket select")
	flag.Parse()

	fmt.Printf("Event Address: %s\n sleeptime:%v", eventAddress, sleetime)
	fmt.Printf("zmq-server-addr: %s\n", server_addr)
	fmt.Printf("chaincodeID: %s\n", chaincodeID)

	socket, err := zmq_wrapper.InitZMQ(server_addr)

	if err != nil {
		fmt.Println("1: ", *err)
		os.Exit(1)
	}

	go systemEventHandler()

	a := createEventClient(eventAddress, listenToRejections, chaincodeID)
	if a == nil {
		fmt.Printf("Error creating event client\n")
		return
	}

	for {
		select {
		case <-a.notfy:
		case r := <-a.rejected:
			zmq_wrapper.SendMsg(socket, rejected_filter, []byte(r.Rejection.ErrorMsg))
		case ce := <-a.cEvent:
			zmq_wrapper.SendMsg(socket, cEvent_filter, ce.ChaincodeEvent.Payload)
		default:
			if break_out {
				fmt.Printf("Safely Exit\n")
				goto End
			}
			time.Sleep(time.Millisecond * 10)
		}
	}
End:
}
