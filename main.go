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
	conn.SetReadDeadline(time.Now().Add(60*time.Second))
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
			fmt.Println("closing")
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
		fmt.Printf("%s\n", message)
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

	handleConnection(conn)



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
