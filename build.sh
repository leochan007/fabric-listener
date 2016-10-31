#!/bin/bash

#LD_LIBRARYPATH=/home/libdev/zmq4/lib 

LD_LIBRARYPATH=/home/libdev/zmq4/lib CGO_CFLAGS=-I/home/libdev/zmq4/include CGO_LDFLAGS=-L/home/libdev/zmq4/lib go build -tags zmq_4_x -o general_listener

#LD_LIBRARYPATH=/home/libdev/zmq4/lib ./ycp_listener -zmq-server-addr=tcp://0.0.0.0:8341 -events-address=127.0.0.1:7053 -listen-to-rejections=true -events-from-chaincode=peersafe-yinchengpai
