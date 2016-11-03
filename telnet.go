package main

import (
	"fmt"
	"net"
	"bufio"
	"strings"
	"errors"
)

func sendTelnetMessage(command string) (string, error) {
	// connect to this socket
	conn, err := net.Dial("tcp", TelnetHost)
	if err != nil {
		fmt.Println("Connection Failure")
		return "", errors.New("Connection Failure")
	}
	reader := bufio.NewReader(conn)
	if err != nil {
		fmt.Println("Reader Failure")
		return "", errors.New("Reader Failure")
	}

	// Authenticate
	line, _, err := reader.ReadLine()
	//fmt.Println(string(line))
	fmt.Fprintf(conn, TelnetPassword + "\n")

	line, _, err = reader.ReadLine()
	//fmt.Println(string(line))
	if strings.HasPrefix(string(line), "Logon successful") {
		fmt.Println("Authentication Failure")
		return "", errors.New("Authentication Failure")
	}

	for i := 0; i < 17; i++ {
		line, _, err = reader.ReadLine()
		//fmt.Println(string(line))
  }

	fmt.Fprintf(conn, "loglevel ALL false\n")
	line, _, err = reader.ReadLine()
	//fmt.Println(string(line))
	line, _, err = reader.ReadLine()
	//fmt.Println(string(line))

	if command == "lp" {
		var lines string
		fmt.Fprintf(conn, command + "\n")
		for {
			line, _, err = reader.ReadLine()
			lineString := string(line)
			if strings.HasPrefix(lineString, "Total of") {
				lines += lineString
				break
			} else {
				lines += strings.Split(lineString, ",")[1] + "\n"
			}
		}

		line = []byte(lines)

	} else {
		fmt.Fprintf(conn, command + "\n")
		line, _, err = reader.ReadLine()
		//fmt.Println(string(line))
	}

	return string(line), nil
}
