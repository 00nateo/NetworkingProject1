package main

import (
	"fmt"
	"os"
	"net"
	"time"
	"strconv"
	"strings"
)

func handleConnection (conn net.Conn){
	conn.SetReadDeadline(time.Now().Add(10*time.Second))
	for {
		buffer := make([]byte, 256)
		n, err := conn.Read(buffer)
		fmt.Printf("Recieved: %s\n", buffer[:n])
		if err != nil{
			fmt.Println("Error reading ", err)
			os.Exit(1)
		}

		msg := strings.TrimSpace(string(buffer[:n]))
		// if msg == "Hello" {
		// 	conn.Write([]byte("Hello to you!\n"))
		// }
		// if !strings.HasSuffix(msg, "\n"){
		// 	conn.Write([]byte(string("no new line error \n")))
		// 	//conn.Close()
		// }
		body := strings.TrimSuffix(msg, "\n")
		segments := strings.Split(body, " ")
		if segments[1] == "BYE"{
			conn.Close()
			break;
		}
		// if len(segments) != 5 || segments[0] != "cs4254fall2026" || segments[1]!= "STATUS" {
		// 	conn.Write([]byte(string("segment errors \n")))
		// 	//conn.Close()
		//
		// }
		ans := 0
		if segments[1] == "STATUS" && isNum(segments[2]) && isNum(segments[4]){
			a, errA := strconv.Atoi(segments[2])
			b, errB := strconv.Atoi(segments[4])
			if errA != nil || errB != nil {
   				conn.Write([]byte("bad number\n"))
   				continue
			}
			// var ans int
			switch segments[3]{
				case "+": ans = a + b
				case "-": ans = a - b
				case "*": ans = a * b
				case "/": ans = a / b
			default: 
				conn.Write([]byte("bad operator\n"))
				continue

			}
			fmt.Printf("answer is %d\n", ans)
		}

		fmt.Println("Sending SOLUTION")
		message := "cs4254fall2026 " + strconv.Itoa(ans) + "\n"
		fmt.Println("%s\n", message)
		_, err2 := conn.Write([]byte(message))
		if err2 != nil{
			fmt.Println("Error", err)
			os.Exit(1)
		}

	}

}

func isNum(s string) bool {
	if s == "" {return false}
	for _, c := range s {
		if c < '0' || c > '9' { return false }
	}
	return true
}
func createConnection(port string , hostname string, pid string){
	// portStr := strconv.Itoa(port)
	host := hostname + ":" + port

	fmt.Printf("Host: %s\n",host)
	fmt.Println("Attempting to connect...")
	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("Error", err)
		os.Exit(1)
	}


	fmt.Println("Sending hello message...")
	message := "cs4254fall2026 HELLO " + pid + "\n"
	_, err2 := conn.Write([]byte(message))
	
	if err2 != nil{
		fmt.Println("Error", err)
		os.Exit(1)
	}
	for{
		handleConnection(conn)
	}


	// fmt.Println("Attempting to read")
	// buffer := make([]byte, 1024)
	// n, err := conn.Read(buffer)
	// if err != nil {
	// 	fmt.Println("Error", err)
	// 	os.Exit(1)
	// }
	//
	// strings := strings.Fields(string(buffer[:n]))
	// for i := 0;i<5;i++{
	// 	fmt.Print( strings[i])
	// 	fmt.Print(" ")
	// }
	// //TODO
	// // calculate()
	//
	// fmt.Printf("\nRecieved from server: %s\n", string(buffer[:n]))


}
func main(){
	//Connect to socket sprinter2.cs.vt.edu
	//Handle server replying with a STATUS message
	//in the STATUS message, extract the maht expression
	//solve math expression
	//expect response of SOLUTION or another STATUS or BYE
	//keep solving expressions until BYE then close the connection
	//once closed it will respond with the secret flag
	//submit code and secret flag
	args := os.Args[1:]
	length := len(args)
	if length != 4{
		fmt.Println("Usage: [-p port] <hostname> <VT Username>")
		os.Exit(1)
	}
	fmt.Println("Args[0]: ", args[0])
	fmt.Println("Args[1]: ", args[1])
	fmt.Println("Args[2]: ", args[2])
	fmt.Println("Args[3]: ", args[3])
	port := args[1]
	hostname := args[2]
	pid := args[3]
	createConnection(port, hostname, pid)
	//Usage: /simpleclient [-p port] <hostname> <VT Username>

	// listener, err := net.Listen("tcp", ":27993")
	// if err != nil {
	// 	fmt.Println("Error listening: ", err)
	// 	return
	// }
	// defer listener.Close()
	// fmt.Println("Server running on :27993")
	// for {
	// 	conn, err := listener.Accept()
	// 	if err != nil{
	// 		fmt.Println("Error accepting:", err)
	// 		continue
	// 	}
	// 	go handleConnection(conn)//one goroutine thread per connection
	// }
//	fmt.Println("Hello world")
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
