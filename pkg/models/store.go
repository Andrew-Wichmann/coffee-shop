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
	OrderCoffee(CoffeeType) (Coffee, error)
	CustomersWaitingToBeServed() (int, error)
	NowServing() (int, error)
	// Sit() error
	// Shit() error
}

type inMemoryStore struct {
	Profit      DollarAmount
	tickets     []Ticket // Should almost certainly be mutexed
	cancel      context.CancelFunc
	wg          sync.WaitGroup
	nowServing  int
	totalServed int
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
	ticket := inMemoryTicket{onCall: onTicketCalled, number: s.totalServed + 1}
	logging.Logger.Debug("Ticket taken", zap.Int("ticket_number", ticket.Number()))
	s.tickets = append(s.tickets, ticket)
	s.totalServed += 1
	return ticket, nil
}

func (s inMemoryStore) OrderCoffee(coffeeType CoffeeType) (Coffee, error) {
	return HouseCoffee{}, nil
}

func (s inMemoryStore) NowServing() (int, error) {
	return s.nowServing, nil
}

func (s *inMemoryStore) Open() {
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if len(s.tickets) > 0 {
					ticket := s.tickets[0]
					s.nowServing = ticket.Number()
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

func (s *inMemoryStore) CustomersWaitingToBeServed() (int, error) {
	return len(s.tickets), nil
}
