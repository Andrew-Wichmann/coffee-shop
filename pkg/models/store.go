package models

import (
	"context"
	"sync"
	"time"

	"github.com/Andrew-Wichmann/coffee-shop/pkg/logging"
	"go.uber.org/zap"
)

func init() {
	CoffeeStore = &inMemoryStore{}
}

var CoffeeStore Store

type Store interface {
	Open()
	Close(context.Context)
	TakeTicket(onTicketCalled func() error) (Ticket, error)
	CheckTicket(Ticket) (int, error)
	OrderCoffee(CoffeeType) (Coffee, error)
	NowServing() (Ticket, error)
	// Sit() error
	// Shit() error
}

type inMemoryStore struct {
	Profit  DollarAmount
	tickets []Ticket // Should almost certainly be mutexed
	cancel  context.CancelFunc
	wg      sync.WaitGroup
	// TODO:
	// tables
	// registers
	// bathrooms
	// dogs
	// blackjack
	// hookers
	// ...
	// forget the tables, registers, bathrooms, and dogs.
}

func (s *inMemoryStore) TakeTicket(onTicketCalled func() error) (Ticket, error) {
	ticket := inMemoryTicket{onCall: onTicketCalled}
	logging.Logger.Debug("Ticket taken", zap.Int("ticket_number", ticket.Number()))
	s.tickets = append(s.tickets, ticket)
	return inMemoryTicket{}, nil
}

func (s inMemoryStore) CheckTicket(ticket Ticket) (int, error) {
	now_serving, err := s.NowServing()
	if err != nil {
		return 0, err
	}
	return ticket.Number() - now_serving.Number(), nil
}

func (s inMemoryStore) OrderCoffee(coffeeType CoffeeType) (Coffee, error) {
	return HouseCoffee{}, nil
}

func (s inMemoryStore) NowServing() (Ticket, error) {
	if len(s.tickets) == 0 {
		return nil, nil
	}
	return s.tickets[0], nil
}

func (s *inMemoryStore) Open() {
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		ticker := time.NewTicker(50 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				logging.Logger.Debug("Checking tickets", zap.Int("ticket_count", len(s.tickets)))
				if len(s.tickets) > 0 {
					ticket := s.tickets[0]
					err := ticket.Call()
					if err != nil {
						logging.Logger.Error("Error calling ticket. Maybe we should consider kicking that customer out?", zap.Error(err))
					}
					s.tickets = s.tickets[1:]
				}
			}
		}
	}()
}

func (s *inMemoryStore) Close(ctx context.Context) {
	logging.Logger.Info("Closing up shop")
	s.cancel()
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()
	select {
	case <-done:
		logging.Logger.Info("Shop closed gracefully")
	case <-ctx.Done():
		logging.Logger.Error("Shop could not close gracefully. Forced to move on.")
	}
}
