package models

import (
	"github.com/Andrew-Wichmann/coffee-shop/pkg/logging"
	"go.uber.org/zap"
)

type Ticket interface {
	Number() int
	Call() error
}

type inMemoryTicket struct {
	onCall func() error
	number int
}

func (ticket inMemoryTicket) Number() int {
	return ticket.number
}

func (ticket inMemoryTicket) Call() error {
	logging.Logger.Info("Now serving", zap.Int("ticket_number", ticket.Number()))
	return ticket.onCall()
}
