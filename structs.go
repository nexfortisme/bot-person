package main

import (
	"fmt"
	"math/rand"
)

// Card holds the card suits and types in the deck
type Card struct {
	Type string
	Suit string
}

func (c Card) String() string {
	return fmt.Sprintf("%s %s", c.Type, c.Suit)
}

// Deck holds the cards in the deck to be shuffled
type Deck []Card

var (
	playing     = true
	deckCounter = 0
)

// New creates a deck of cards to be used
func New() (deck Deck) {

	// Valid types include Two, Three, Four, Five, Six
	// Seven, Eight, Nine, Ten, Jack, Queen, King & Ace
	types := []string{"Two", "Three", "Four", "Five", "Six", "Seven",
		"Eight", "Nine", "Ten", "Jack", "Queen", "King", "Ace"}

	// Valid suits include Heart, Diamond, Club & Spade
	suits := []string{"Heart", "Diamond", "Club", "Spade"}

	// Loop over each type and suit appending to the deck
	for i := 0; i < len(types); i++ {
		for n := 0; n < len(suits); n++ {
			card := Card{
				Type: types[i],
				Suit: suits[n],
			}
			deck = append(deck, card)
		}
	}
	return
}

// Shuffle the deck
func Shuffle(d Deck) Deck {
	for i := 1; i < len(d); i++ {
		// Create a random int up to the number of cards
		r := rand.Intn(i + 1)

		// If the the current card doesn't match the random
		// int we generated then we'll switch them out
		if i != r {
			d[r], d[i] = d[i], d[r]
		}
	}
	return d
}

// Deal a specified amount of cards
func Deal(d Deck, n int) []Card {

	var hand []Card

	for i := deckCounter; i < n+deckCounter; i++ {
		// fmt.Println(d[i])
		hand = append(hand, d[i])
	}

	deckCounter += n

	return hand
}