package smtp

import (
	"log/slog"
	"net"
	"net/textproto"
	"strings"
)

type server struct {
	config     ServerConfig
	clientName string
}

type ServerConfig struct {
	Hostname string
}

func CreateServer(config ServerConfig) server {
	return server{config: config}
}

func (s *server) HandleConnection(conn net.Conn) {
	defer conn.Close()
	text := textproto.NewConn(conn)
	defer text.Close()
	s.HandleRequest(text)
}

func (s *server) HandleRequest(text *textproto.Conn) {
	err := text.PrintfLine("220 OK")
	if err != nil {
		slog.Error("error sending result", "msg", err)
	}

	for {
		line, err := text.ReadLine()
		if err != nil {
			slog.Error("error reading line", "msg", err)
		}
		slog.Info(line)
		s.handleCommand(line)
	}
}

func (s *server) handleCommand(command string) {
	idx := strings.Index(command, " ")
	if idx == -1 {
		slog.Error("error getting command from line", "command", command)
		return
	}

	cmd := command[0:idx]
	switch cmd {
	case "EHLO":
		s.handleEhlo(command[idx+1:])
	default:
		slog.Error("no handler for command ", "cmd", cmd)
	}
}

func (s *server) handleEhlo(data string) {
	s.clientName = data
}
