package smtp

import (
	"log/slog"
	"net"
	"net/textproto"
)

func HandleRequest(conn net.Conn) {
	defer conn.Close()
	text := textproto.NewConn(conn)
	defer text.Close()
	err := text.PrintfLine("220 OK")
	if err != nil {
		slog.Error("error", err)
	}
}
