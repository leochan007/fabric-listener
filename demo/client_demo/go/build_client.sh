#!/bin/sh

CGO_CFLAGS=-I/home/libdev/zmq4/include CGO_LDFLAGS=-L/home/libdev/zmq4/lib go build -tags zmq_4_x -o gozmq_client
./gozmq_client -zmq-server-addr=tcp://127.0.0.1:8341