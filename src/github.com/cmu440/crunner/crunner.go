package main

import (
	"bufio"
    "fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
    defaultHost = "localhost"
    defaultPort = 9999
)

// To test your server implementation, you might find it helpful to implement a
// simple 'client runner' program. The program could be very simple, as long as
// it is able to connect with and send messages to your server and is able to
// read and print out the server's response to standard output. Whether or
// not you add any code to this file will not affect your grade.
func main() {
	addr := strings.Join([]string{defaultHost, strconv.Itoa(defaultPort)}, ":")
	conn, err := net.Dial("tcp", addr)

	if err != nil {
		log.Fatalln(err)
	}


	defer conn.Close()
	timeoutDuration := 3 * time.Second
	conn.SetReadDeadline(time.Now().Add(timeoutDuration))

	buff := make([]byte, 1024)
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("Enter the message >> ")
		text, _ := reader.ReadString('\n')
		conn.Write([]byte(text))
		log.Printf("Send: %s\n", text)
		n, err := conn.Read(buff)
		if nil != err {
			fmt.Printf("Received: %s\n\n", buff[:n])
		}
	}
	// scanner := bufio.NewScanner(os.Stdin)
	// fmt.Printf("Enter the message >> ")
	// for scanner.Scan() {
	// 	conn.Write([]byte(scanner.Bytes()))
	// 	fmt.Printf("Sent: %s\n", scanner.Text())
	// 	fmt.Printf("Enter the message >> ")
	// }
}
