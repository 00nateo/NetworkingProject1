package main

import (
	"fmt"
	"net"
	"time"

)

func handleConnection (conn net.Conn){
	defer conn.Close()
	conn.SetReadDeadline(time.Now().Add(10*time.Second))
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil{
			fmt.Println("Error reading ", err)
			return
		}
		fmt.Printf("Recieved:: %s", buffer[:n])
		conn.Write([]byte("Message recieved\n"))
	}

}
func main(){
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error listening: ", err)
		return
	}
	defer listener.Close()
	fmt.Println("Server running on :8080")
	for {
		conn, err := listener.Accept()
		if err != nil{
			fmt.Println("Error accepting:", err)
			continue
		}
		go handleConnection(conn)//one goroutine thread per connection
	}
	fmt.Println("Hello world")
}
/*Once the socket is connected, the client sends a HELLO message to the
server. The format of the HELLO message is:
cs4254fall2026 HELLO [your VT username]\n
In your program you should replace [your VT username] with your actual VT Username. You
must supply your VT Username so the server can look up the appropriate secret flag for you. The
server will reply with a STATUS message. The format of the STATUS message is:
1
cs4254fall2026 STATUS [a number] [a math operator] [another number]\n
The three variable fields represent a simple mathematical expression, e.g., "5 + 10". The server
may return plus, minus, multiplication, or division expressions. All numbers will be between 1 and
1000. Your program must solve the mathematical expression and return the answer to the server
in a SOLUTION message. The SOLUTION message has the following format:
cs4254fall2026 [the solution]\n
It is okay for the solution to be negative. In the case of division, round the answer down to the
nearest integer (do not send floating point numbers to the server).
The server will respond to the SOLUTION message with either another STATUS message, or
a BYE message. If the server terminates the connection, that means your solution was incorrect.
If the server sends another STATUS message, your program must solve the expression and return
another SOLUTION message. The server will ask your program to solve hundreds of expressions;
the exact number of expressions is chosen at random. Eventually, the server will return a BYE
message. The BYE message has the following format:
cs4254fall2026 [a 64 byte secret flag] BYE\n
Once your program has received the BYE message, it can close the connection to the server. If
the server returns "Unknown_VT_Username" in the BYE message, that means it did not recognize
the VT Username that you supplied in the HELLO message. Otherwise, the 64-byte string is your
secret flag: write this value down, since you need to turn it in along with your code.
*/
