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
	text       *textproto.Conn
}

type ServerConfig struct {
	Hostname string
}

func CreateServer(config ServerConfig) server {
	return server{config: config}
}

func (s *server) HandleConnection(conn net.Conn) {
	s.text = textproto.NewConn(conn)

	err := s.text.PrintfLine("220 OK")
	if err != nil {
		slog.Error("error sending connection result", "msg", err)
	}

	s.HandleRequest()
}

func (s *server) HandleRequest() {
	for {
		line, err := s.text.ReadLine()
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

	err := s.text.PrintfLine("250-%v greets %v", s.config.Hostname, s.clientName)
	if err != nil {
		slog.Error("error sending greeting", "err", err)
	}

	err = s.text.PrintfLine("250-8BITMIME")
	if err != nil {
		slog.Error("error sending greeting extension info", "err", err)
	}

	err = s.text.PrintfLine("250 HELP")
	if err != nil {
		slog.Error("error sending greeting end", "err", err)
	}
}
