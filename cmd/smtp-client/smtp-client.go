package main

import (
	"log/slog"
	"net/smtp"
)

func main() {
	client, err := smtp.Dial("localhost:8080")
	if err != nil {
		slog.Error("error dialing", err)
	}
	err = client.Noop()
	if err != nil {
		slog.Error("error noop", err)
	}
	err = client.Mail("test@test.test")
	if err != nil {
		slog.Error("error mail", err)
	}
}
