package main

import (
	"GolangPractice/hw_12/internal/const"
	"GolangPractice/hw_12/internal/llogger"
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"os"
)

var logger = llogger.SetupLogger()

func main() {
	us := bufio.NewScanner(os.Stdin)
	uw := bufio.NewWriter(os.Stdout)

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(fmt.Errorf("error connecting: %w", err))
	}
	logger.Info("Client connected to server")
	defer conn.Close()

	sr := bufio.NewReader(conn)
	sw := bufio.NewWriter(conn)

	uw.WriteString(_const.HelpText)
	uw.Flush()
	fmt.Print(">>> ")
	for us.Scan() {
		msg := us.Text()

		fmt.Printf("User Input: %s\n", msg)

		sw.WriteString(base64.StdEncoding.EncodeToString([]byte(msg)) + "\n")
		sw.Flush()

		resp, err := sr.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				logger.Info("Server closed the connection (EOF)")
				break
			} else {
				logger.Error(fmt.Sprintf("error reading response: %v", err))
				break
			}
		}

		uw.WriteString(resp)
		uw.Flush()
		fmt.Print(">>> ")
	}
}
