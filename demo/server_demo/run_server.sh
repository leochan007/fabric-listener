#!/bin/sh

CGO_CFLAGS=-I/home/libdev/zmq4/include CGO_LDFLAGS=-L/home/libdev/zmq4/lib go run -tags zmq_4_x demo.go