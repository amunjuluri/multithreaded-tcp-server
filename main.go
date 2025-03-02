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
	fmt.Println("Processing the request")
	time.Sleep(5 * time.Second) // sleeps for 5 seconds and then starts writing response

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
//  blocking mean, suppose when your client is not reading while you are writing,
// your process is blocked

func main() {
	listener, err := net.Listen("tcp", ":1729") //? 1. reserving a port, now we need a client to connect
	if err != nil {
		log.Fatal(err)
	}

	for { //? 5. so, to continuosly accept requests, we put this in an infinite for loop
		fmt.Println("Waiting for Client to Connect")
		conn, err := listener.Accept() //? 2. now it waits for someone to connect, now we need to read response, write the response, and close the connection
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Client Connected")

		go do(conn) //? 7 create a simple go routine and it would run a seperate go routine

	}
}

//? 6. but now, this is a single threaded server, it cannot handle concurrent request, it can only handle one by one
//? what we can do is, while we accept the connection, and processing it with do function
//? what if we do that process, the function "do" part, in a seperate thread
//? so that the main thread, again goes to the top while it processes and again waits to accept the client

//!! problems:
// - when large number of clients connect, there could be a chance of thread overload
//? optimzations:
// - don't spin up a new thread for every request, limit maximum number of threads
// - add thread pool(lets say 1000 threads,:= go to thread pool -> process with a thread -> put back the thread in thread pool), to save thread creation time
// - add a timeout - what if a client is connected and the client never sends you a request, so you need to kill that connection after certain amount of time
// - tcp backlog queue - an OS level setting that determines how many connections can be holded in the queue
/*
What are these queued connections?

In simple words, the backlog parameter specifies the number of pending connections the queue will hold.

When multiple clients connect to the server, the server then holds the incoming requests in a queue. The clients are arranged in the queue, and the server processes their requests one by one as and when queue-member proceeds. The nature of this kind of connection is called queued connection.

Does it make any difference for client requests? (I mean is the server that is running with socket.listen(5) different from the server that is running with socket.listen(1) in accepting connection requests or in receiving data?)

Yes, both cases are different. The first case would allow only 5 clients to be arranged to the queue; whereas in the case of backlog=1, only 1 connection can be hold in the queue, thereby resulting in the dropping of the further connection request!
*/
