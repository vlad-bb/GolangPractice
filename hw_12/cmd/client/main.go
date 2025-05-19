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
	"strings"
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

		elems := strings.SplitN(msg, " ", 2)
		cmd := elems[0]

		var payload string
		if len(elems) > 1 {
			// Тільки payload кодуємо в base64
			payload = base64.StdEncoding.EncodeToString([]byte(elems[1]))
			sw.WriteString(fmt.Sprintf("%s %s\n", cmd, payload))
		} else {
			// Команди без payload (наприклад, list, exit)
			sw.WriteString(fmt.Sprintf("%s\n", cmd))
		}

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
