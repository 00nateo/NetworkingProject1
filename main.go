package main

import (
	"fmt"
	"os"
	"net"
	"strconv"
	"strings"
)

func handleConnection (conn net.Conn){
	// conn.SetReadDeadline(time.Now().Add(60*time.Second))
	for {
		//create buffer and read connection from connection
		buffer := make([]byte, 256)
		n, err := conn.Read(buffer)
		fmt.Printf("Recieved: %s", buffer[:n])
		if err != nil{
			fmt.Println("Error reading ", err)
			os.Exit(1)
		}

		//Parse the message
		msg := strings.TrimSpace(string(buffer[:n]))
		body := strings.TrimSuffix(msg, "\n")
		segments := strings.Split(body, " ")

		//check for BYE message (last message)
		if len(segments) == 3 && segments[2] == "BYE"{
			if (segments[0] != "cs4254fall2026"){
				fmt.Fprintln(os.Stderr, "Error: Unknown Message Type")
				os.Exit(1)
			}
			if (segments[1] == "Unknown_VT_Username"){
				fmt.Fprintln(os.Stderr, "Error: Unknown Vt User")
				os.Exit(1)
			}
			fmt.Printf("Final secret flag: %s\n", segments[1])
			conn.Close()
			break;
		}

		ans := 0
		// check STATUS and parse arguments to complete calculation
		if segments[1] == "STATUS" && 
				isNum(segments[2]) &&
				isNum(segments[4]) &&
				segments[0] == "cs4254fall2026" {
			a, errA := strconv.Atoi(segments[2])
			b, errB := strconv.Atoi(segments[4])
			// handle strconv.Atoi return possibilities (num or err)
			if errA != nil || errB != nil {
   				conn.Write([]byte("NAN in message\n"))
   				continue
			}
			// do specified calculation
			switch segments[3]{
				case "+": ans = a + b
				case "-": ans = a - b
				case "*": ans = a * b
				case "/": ans = a / b
			default: 
				conn.Write([]byte("Operator not one of +, -, *, /\n"))
				continue
			}
			fmt.Printf("Calculated answer is %d\n", ans)
		}
		//print ans
		message := "cs4254fall2026 " + strconv.Itoa(ans) + "\n"
		fmt.Printf("Sending: %s\n", message)
		_, err2 := conn.Write([]byte(message))
		if err2 != nil{
			fmt.Println("Error", err)
			os.Exit(1)
		}

	}

}
// return true if given string is comprised of integers between 0 and 9, ie 123 is 1, 2 and 3
func isNum(s string) bool {
	if s == "" {return false}
	for _, c := range s {
		if c < '0' || c > '9' { return false }
	}
	return true
}

//creates initial TCP connection and sends HELLO message
//Takes in port number, hostname, and VT pid
func createConnection(port string , hostname string, pid string){
	host := hostname + ":" + port
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

//Connect to socket sprinter(2/3).cs.vt.edu
//Handle server replying with a STATUS message
//in the STATUS message, extract the maht expression
//solve math expression
//expect response of SOLUTION or another STATUS or BYE
//keep solving expressions until BYE then close the connection
//once closed it will respond with the secret flag
//submit code and secret flag
func main(){
	args := os.Args[1:]
	length := len(args)
	if length != 4 && length != 2{
		fmt.Println("Usage: [-p (optional) port] <hostname> <VT Username>")
		os.Exit(1)
	}

	//optional -p port param. Default port: 27993
	if (length == 4 && args[0] == "-p"){
		port := args[1]
		hostname := args[2]
		pid := args[3]
		createConnection(port, hostname, pid)
	} else if (length == 2){
		port :="27993"
		hostname := args[0]
		pid := args[1]
		createConnection(port, hostname, pid)
	} else {
		os.Exit(1)
	}
}
