package main

import (
	"fmt"
	"net"
	"strings"
	"time"

	curr "github.com/vladimirvivien/learning-go/ch11/curr0"
)

var currencies = curr.Load("data.csv")

func main() {
	ln, _ := net.Listen("tcp", ":4040")
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			conn.Close()
			return
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		cmdLine := make([]byte, (1024 * 4))
		n, err := conn.Read(cmdLine)
		if n == 0 || err != nil {
			return
		}
		cmd, param := parseCommand(string(cmdLine[0:n]))
		if cmd == "" {
			continue
		}

		switch strings.ToUpper(cmd) {
		case "GET":
			result := curr.Find(currencies, param)

			for _, cur := range result {
				_, err := fmt.Fprintf(
					conn,
					"%s %s %s %s\n",
					cur.Name, cur.Code,
					cur.Number, cur.Country,
				)
				if err != nil {
					return
				}

				conn.SetWriteDeadline(
					time.Now().Add(time.Second * 5))
			}
			conn.SetWriteDeadline(
				time.Now().Add(time.Second * 300))

		default:
			conn.Write([]byte("Invalid Command\n"))
		}

	}
}

func parseCommand(cmdLine string) (cmd, param string) {
	parts := strings.Split(cmdLine, " ")
	if len(parts) != 2 {
		return "", ""
	}
	cmd = strings.TrimSpace(parts[0])
	param = strings.TrimSpace(parts[1])
	return
}
