package main

import (
	_const "GolangPractice/hw_13/internal/const"
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	cmds "GolangPractice/hw_13/internal/commands"
	"GolangPractice/hw_13/internal/document_store"
	"GolangPractice/hw_13/internal/like_mongo"
	"GolangPractice/hw_13/internal/llogger"
	"github.com/brianvoe/gofakeit/v6"
)

var logger = llogger.SetupLogger()

func execPut(raw string, rs like_mongo.Service) (string, error) {
	p := &cmds.PutCommandRequestPayload{}
	err := json.Unmarshal([]byte(raw), p)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling payload: %w", err)
	}

	_, err = rs.CreateRecord(p.Key, p.Value)
	if err != nil {
		return "", err
	}

	resp := &cmds.PutCommandResponsePayload{}
	rawResp, err := json.Marshal(resp)
	if err != nil {
		return "", fmt.Errorf("error marshalling response: %w", err)
	}

	return string(rawResp), nil
}

func execGet(raw string, rs like_mongo.Service) (string, error) {
	p := &cmds.GetCommandRequestPayload{}
	err := json.Unmarshal([]byte(raw), p)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling payload: %w", err)
	}

	record, err := rs.GetRecord(p.Key)
	if err != nil {
		return "", err
	}

	resp := &cmds.GetCommandResponsePayload{
		Value: record.Value,
		Ok:    true,
	}

	rawResp, err := json.Marshal(resp)
	if err != nil {
		return "", fmt.Errorf("error marshalling response: %w", err)
	}

	return string(rawResp), nil
}

func execDelete(raw string, rs like_mongo.Service) (string, error) {
	p := &cmds.DeleteCommandRequestPayload{}
	err := json.Unmarshal([]byte(raw), p)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling payload: %w", err)
	}

	err = rs.DeleteRecord(p.Key)
	if err != nil {
		return "", err
	}
	resp := &cmds.DeleteCommandResponsePayload{
		Ok: true,
	}

	rawResp, err := json.Marshal(resp)
	if err != nil {
		return "", fmt.Errorf("error marshalling response: %w", err)
	}

	return string(rawResp), nil
}

func execInvalidCommand(raw string) (string, error) {
	resp := &cmds.ExitCommandResponsePayload{
		Msg: raw,
	}
	rawResp, err := json.Marshal(resp)
	if err != nil {
		return "", fmt.Errorf("error marshalling response: %w", err)
	}

	return string(rawResp), nil
}

func execList(rs like_mongo.Service) (string, error) {
	records, err := rs.ListRecords()
	if err != nil {
		return "", err
	}
	resp := &cmds.ListCommandResponsePayload{
		Items: records,
	}
	rawResp, err := json.Marshal(resp)
	if err != nil {
		return "", fmt.Errorf("error marshalling response: %w", err)
	}

	return string(rawResp), nil
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	logger.Info("connection accepted" + conn.RemoteAddr().String())
	scanner := bufio.NewScanner(conn)
	w := bufio.NewWriter(conn)

	var store = document_store.NewStore()

	gofakeit.Seed(0)
	name := gofakeit.Word()
	var rs = like_mongo.CreateService(store, "key", name)

	for scanner.Scan() {
		msg := scanner.Text()

		req, _ := base64.StdEncoding.DecodeString(msg)
		elems := strings.Split(string(req), " ")

		if string(req) == _const.Exit {
			w.WriteString(_const.ExitMsg)
			w.Flush()
			return
		}

		if len(elems) > 2 {
			w.WriteString("invalid command\n")
			w.Flush()
			continue
		}

		var resp string
		var err error

		switch elems[0] {
		case cmds.PutCommandName:
			resp, err = execPut(elems[1], rs)
		case cmds.GetCommandName:
			resp, err = execGet(elems[1], rs)
		case cmds.DeleteCommandName:
			resp, err = execDelete(elems[1], rs)
		case cmds.ListCommandName:
			resp, err = execList(rs)
		default:
			resp, err = execInvalidCommand("invalid command\n" + _const.HelpText)
		}

		if err != nil {
			w.WriteString(fmt.Sprintf("error: %s\n", err))
		}

		w.WriteString(fmt.Sprintf("response: %s\n", resp))

		w.Flush()
	}

	fmt.Println("connection closed")
}

func main() {
	addr, shutdown, err := StartTCPServer("127.0.0.1:8080")
	logger.Info("Server started on %s", addr)
	if err != nil {
		logger.Error("Failed to start server: %v", err)
	}
	defer shutdown()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	logger.Info("Shutdown signal received")
	shutdown()
	logger.Info("Server gracefully stopped")
}

func StartTCPServer(addr string) (net.Addr, func(), error) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, nil, fmt.Errorf("error listening: %w", err)
	}

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				return
			}
			go handleConnection(conn)
		}
	}()

	return l.Addr(), func() { l.Close() }, nil
}
