package models

type CoffeeType int

const (
	HOUSE CoffeeType = iota
	CORTADO
)

type Coffee interface {
	Type() CoffeeType
	Name() string
	Price() DollarAmount
}

type HouseCoffee struct{}

func (HouseCoffee) Type() CoffeeType {
	return HOUSE
}

func (HouseCoffee) Name() string {
	return "House black coffee"
}

func (HouseCoffee) Price() DollarAmount {
	return DollarAmount{
		Dollars: 1,
		Cents:   50,
	}
}
