package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func do(conn net.Conn) {
	//? 3. to read, we need to read the req, and store it somewhere, so,
	buff := make([]byte, 1024) //make a byte array,
	_, err := conn.Read(buff)  // this waits, until the clients sends the request, returns int(number of bytes read), and error
	if err != nil {
		log.Fatal(err)

	}
	time.Sleep(5 * time.Second ) // sleeps for 5 seconds and then starts writing response

	conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\nHello World\r\n"))
	//The \r\n sequences are crucial parts of the HTTP protocol specification:
	// 	proper formatting is required for browsers and HTTP clients to correctly interpret the message:
	// The status line must end with \r\n
	// Each header line must end with \r\n
	// Headers must be separated from the body with an empty line (which is \r\n\r\n)
	// The body content can also include line breaks

	conn.Close() //close the connection
}

//? 4. read and write are blocking calls,
// and we need to be sure when call these functions,
//  when your client is not reading while you are writing,
// your process is blocked

func main() {
	listener, err := net.Listen("tcp", ":1729") //? 1. reserving a port, now we need a client to connect
	if err != nil {
		log.Fatal(err)
	}

	for { //? 5. so, to continuosly accept requests, we put this in an infinite for loop
		conn, err := listener.Accept() //? 2. now it waits for someone to connect, now we need to read response, write the response, and close the connection
		if err != nil {
			log.Fatal(err)
		}

		do(conn)
		fmt.Println("zuiii")
	}
}
