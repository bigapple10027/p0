// Implementation of a KeyValueServer. Students should write their code in this file.

package p0

import (
	"fmt"
	"log"
)
type keyValueServer struct {
    // TODO: implement this!
	connectionChannel chan Conn
}

// New creates and returns (but does not start) a new KeyValueServer.
func New() KeyValueServer {
	var ret = keyValueServer{}
	// TODO: initialization here
	return ret
}

func (kvs *keyValueServer) Start(port int) error {
    // TODO: go handler that handles all the request with channels.
	listener, err := net.Listen("tcp", ":" + iota(port))
	if err != nil {
		fmt.Println("keyValueServre.Start: Can't create a listener, error occured")
		return err
	}
	kvs.listener = listener

	go runServer(kvs)
    return nil
}

func (kvs *keyValueServer) Close() {
    // TODO: implement this!
}

func (kvs *keyValueServer) Count() int {
    // TODO: implement this!
    return -1
}

// TODO: add additional methods/functions below!

func runServer(kvs *keyValueServer) {
	for {
		select {
			// Handle when there is new connection.
			case x < x:
				// Add the client to kvs.
				// 

			// Handle when client wants to access the underlying database
			case 

			//  Handle when client tries to disconnect
			case

			// Handle when client 		
		}
	}
}

func listenToConnections(listener *Listener) {
	for {
		connection, err := listener.Accept()
		if nil != err {
			// handle err
		} else {
			// notify server handler thread that there is a new connection.
			// Pass the new connection back through channel
		}
	}
}
