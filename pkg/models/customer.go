package models

type Customer struct {
	Name    string
	Coffees []Coffee
	Mood    CustomerMood
}
type CustomerMood int

const (
	Happy CustomerMood = iota
	Sad
	Neutral
	Angry
	Flirtatious
)
