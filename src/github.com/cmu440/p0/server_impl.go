// Implementation of a KeyValueServer. Students should write their code in this file.

package p0

import (
	"bufio"
	"bytes"
	"fmt"
	// "log"
	"io"
	"net"
	"strconv"
	"strings"
)

const (
	MAX_NUM_BYTES = 1024
	MAX_MESSAGE_QUEUE_LENGTH = 500
)

type DBRequest struct {
	isGet bool
	key string
	value []byte
}

// type DBRequest struct {
// 	isGet bool
// 	key string
// 	value string
// }

type CountClientsRequest struct {
	countChannel chan int
}

type Client struct {
	connection net.Conn
	bytesChannel chan []byte
	quitReadChannel chan bool
	quitWriteChannel chan bool
}

type keyValueServer struct {
    // TODO: implement this!
	connectionChannel chan net.Conn
	dbRequestChannel chan *DBRequest
	countClientsRequestChannel chan CountClientsRequest
	deadClient chan *Client
	quitServer chan bool
	quitListening chan bool
	clients []*Client
	listener net.Listener
}

// New creates and returns (but does not start) a new KeyValueServer.
func New() KeyValueServer {
	var ret = &keyValueServer{}
	ret.connectionChannel = make(chan net.Conn)
	ret.dbRequestChannel = make(chan *DBRequest)
	ret.countClientsRequestChannel = make(chan CountClientsRequest)
	ret.quitListening = make(chan bool)
	ret.quitServer = make(chan bool)
	ret.clients = make([]*Client, 0, 10)

	return ret
}

func (kvs *keyValueServer) Start(port int) error {
    // TODO: go handler that handles all the request with channels.
	listener, err := net.Listen("tcp", ":" + strconv.Itoa(port))
	fmt.Printf("Start server on port: %d\n", port)
	if err != nil {
		fmt.Println("keyValueServre.Start: Can't create a listener, error occured")
		return err
	}
	kvs.listener = listener
	init_db()

	go runServer(kvs)
	go listenToConnections(kvs)
    return nil
}

func (kvs *keyValueServer) Close() {
	fmt.Println("Close() got called")
	kvs.listener.Close()
	kvs.quitServer <- true
	kvs.quitListening <- true
}

func (kvs *keyValueServer) Count() int {
	countClientsRequest := CountClientsRequest{make(chan int)}
	kvs.countClientsRequestChannel <- countClientsRequest
    return <-countClientsRequest.countChannel
}

/*******************************************/
// Private keyValueServer member functions
/*******************************************/


/********************************************/
// End(Private keyValueSever member functions)
/********************************************/

// TODO: add additional methods/functions below!
func runServer(kvs *keyValueServer) {
	for {
		select {

		// Handle when there is a new connection.
		case newConnection := <-kvs.connectionChannel:
			// fmt.Println("Server got new connection ...")
			client := &Client{
							connection: newConnection,
							bytesChannel: make(chan []byte, MAX_MESSAGE_QUEUE_LENGTH),
							quitReadChannel: make(chan bool),
							quitWriteChannel: make(chan bool)}
			kvs.clients = append(kvs.clients, client)
			go read(kvs, client)
			go write(client)

		case deadClient := <-kvs.deadClient:
			fmt.Println("Deadclient signaled!!")
			for i, client := range kvs.clients {
				if client == deadClient {
					kvs.clients = append(kvs.clients[:i], kvs.clients[i+1:]...)
					fmt.Println("Removed one deadclient")
					break
				}
			}

		// Handle when client wants to access the underlying database
		case dbRequest := <-kvs.dbRequestChannel:
			if dbRequest.isGet {
				valueBytes := get(dbRequest.key)
				// fmt.Printf("key: %s\n", dbRequest.key)
				// fmt.Printf("value: %s\n", string(valueBytes[:]))

				outputBytes := []byte(strings.Join([]string{dbRequest.key, string(valueBytes[:])}, ","))
				for _, client := range(kvs.clients) {
					// fmt.Printf("Write <%v> to client\n", string(outputBytes))
					client.bytesChannel <- outputBytes
				}
			} else {
				put(dbRequest.key, []byte(dbRequest.value))
			}

		//  Handle when client tries to disconnect
		case countClientsRequest := <-kvs.countClientsRequestChannel:
			countClientsRequest.countChannel <- len(kvs.clients)

		case <-kvs.quitServer:
			fmt.Println("Quiting server...")
			for _, client := range(kvs.clients) {
				client.connection.Close()
				// fmt.Println("Sending signals to quit reading and writing channels")
				client.quitReadChannel <- true
				client.quitWriteChannel <- true
				// fmt.Println("Read and write signals sent successfully")
			}
			fmt.Println("Quited the server")
			return


		// Handle when client disconnect
		}
	}
}

func listenToConnections(kvs *keyValueServer) {
	for {
		select {
		case <-kvs.quitListening:
			fmt.Println("Quit listenting")
			return
			
		default:
			connection, err := kvs.listener.Accept()
			if nil == err {
				// notify server handler thread that there is a new connection.
				// Pass the new connection back through channel.
				kvs.connectionChannel <- connection
				// fmt.Println("Got new connections")
			}
		}
	}
}

func read(kvs *keyValueServer, client *Client) {
	// create reader from client.connection
	bufReader := bufio.NewReader(client.connection)
	for {
		select {
		case <-client.quitReadChannel:
			fmt.Println("Quiting the read thread")
			return
		default:
			buffer, err := bufReader.ReadBytes('\n')

			if err == io.EOF {
				kvs.deadClient <- client
				fmt.Println("Found deadClient, signaled!!")
			} else if nil != err {
				fmt.Printf("Reading error: %v", err)
				return
			} else  {

				dbRequest := createDbRequest(buffer)
				kvs.dbRequestChannel <- &dbRequest
			}
		}
	}
}

func write(client *Client) {
	for {
		select {
		case <-client.quitWriteChannel:
			fmt.Println("Quiting the write thread")
			return
		case bytes := <-client.bytesChannel:
			client.connection.Write(bytes)
		}
	}
}

func createDbRequest(message []byte) DBRequest {
	tokens := bytes.Split(message, []byte(","))

	if string(tokens[0]) == "put" {
		key := string(tokens[1][:])

		// do a "put" query
		return DBRequest{
			isGet: false,
			key:   key,
			value: tokens[2],
		}
	} else {
		// remove trailing \n from get,key\n request
		keyBin := tokens[1][:len(tokens[1])-1]
		key := string(keyBin[:])
		return DBRequest {
			isGet: true,
			key: key,
		}
	}
}

// func createDbRequest(bytes []byte) DBRequest {
// 	request := string(bytes[:])
// 	fmt.Printf("Database request: %v\n", request)
// 	tokens := strings.Split(request, ",")
// 	// fmt.Printf("Tokens: %v \n", tokens)
// 	
// 	dbRequest := DBRequest{}
// 	if tokens[0] == "get" {
// 		dbRequest.isGet = true
// 		dbRequest.key = strings.TrimSuffix(tokens[1], "\n")
// 	} else {
// 		dbRequest.isGet = false
// 		dbRequest.key = tokens[1]
// 		dbRequest.value = tokens[2]
// 	}
// 	return dbRequest
// }
