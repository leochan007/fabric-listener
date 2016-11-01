#!/bin/bash

ver=4.1.5

sudo apt-get install -y pkg-config

if [ ! -d "./zeromq-$ver" ]; then
echo 'dir not existed.'
wget --no-check-certificate https://github.com/zeromq/zeromq4-1/releases/download/v$ver/zeromq-$ver.tar.gz
tar -xzvf zeromq-$ver.tar.gz
cd zeromq-$ver
./configure --prefix=/home/libdev/zmq4
make
else
echo 'dir existed. entering...'
cd zeromq-$ver
make clean
make
fi

sudo make install
#sudo sh -c 'echo "/home/libdev/zmq4" >> /etc/ld.so.conf'
#sudo ldconfig
sudo ln -s /home/libdev/zmq4/lib/libzmq.so /usr/lib/libzmq.so
export PKG_CONFIG_PATH=`pwd`/zeromq-4.1.5/src
#CGO_CFLAGS=-I/home/libdev/zmq4/include CGO_LDFLAGS=-L/home/libdev/zmq4/lib 
go get -tags zmq_4_x github.com/alecthomas/gozmq

