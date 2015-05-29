package main

import (
	//"errors"
	"fmt"
	"math/rand"
	//"os"
	"runtime"
	//"strconv"
	"sync"
	"time"
)

var wg sync.WaitGroup
var roundsPlayed = make(chan []Card)

func main() {
	//Parse command line arguments
	//size, trials, err := parseArgs()
	runtime.GOMAXPROCS(6)
	rand.Seed(time.Now().Unix())
	//Monte Carlo Everything
	size := 10
	trials := 30
	northCards := []Card{Card{"ah"}, Card{"kh"}, Card{"qh"}, Card{"jh"}} //, Card{"th"}}
	eastCards := []Card{Card{"ac"}, Card{"kc"}, Card{"qc"}, Card{"jc"}}  //, Card{"tc"}}
	southCards := []Card{Card{"ad"}, Card{"kd"}, Card{"qd"}, Card{"jd"}} //, Card{"td"}}
	westCards := []Card{Card{"as"}, Card{"ks"}, Card{"qs"}, Card{"js"}}  //, Card{"ts"}}
	players := []Player{Player{"north", northCards}, Player{"east", eastCards}, Player{"south", southCards}, Player{"west", westCards}}
	probSorted(size, trials, players)
	return
}

type Player struct {
	Name string
	Hand []Card
}

func (p *Player) CopyOf() Player {
	t := make([]Card, len(p.Hand))
	copy(t, p.Hand)
	return Player{p.Name, t}
}

func (p *Player) removeCard(index int) Card {
	result := p.Hand[index]
	p.Hand = append(p.Hand[:index], p.Hand[index+1:]...)
	return result
}

type Card struct {
	Name string
}

type Round struct {
	Players         []Player
	CurrentPosition int
	CardsPlayed     []Card
}

func (r *Round) playNextCard(index int) {
	cardToPlay := r.Players[r.CurrentPosition].removeCard(index)
	r.CardsPlayed = append(r.CardsPlayed, cardToPlay)
	r.CurrentPosition = (r.CurrentPosition + 1) % len(r.Players)
	//fmt.Println("cards played are %v", r.Players[0])
}

func (r *Round) Done() bool {
	lastPlayerIndex := len(r.Players) - 1
	lastPlayerCardsLeft := len(r.Players[lastPlayerIndex].Hand)
	return (lastPlayerCardsLeft == 0)
}

func (r *Round) CopyOf() Round {
	//b := append([]Player(nil), r.Players)
	b := make([]Player, 0)
	for _, player := range r.Players {
		b = append(b, player.CopyOf())
	}

	c := make([]Card, len(r.CardsPlayed))
	copy(c, r.CardsPlayed)

	return Round{b, r.CurrentPosition, c}
}

func probSorted(size, trials int, players []Player) {
	//Make channel to relay whether deck is sorted
	//sorted := make(chan bool)
	//Number of times deck shows up true
	//truth := 0
	combinations := 0
	//for i := 0; i < trials; i++ {
	//Increments waitgroup
	startingRound := Round{players, 0, make([]Card, 0)}
	//roundCopy := startingRound.CopyOf()
	//fmt.Printf("roundCopy after has these hands :: %v%v\n", roundCopy.CardsPlayed, roundCopy.Players)
	//fmt.Printf("startingRound after has these hands :: %v%v\n", startingRound.CardsPlayed, startingRound.Players)

	//wg.Add(1)
	go playRound(&startingRound)
	// wg.Add(1)
	// go playRound(&roundCopy)
	//shuffle a new virtual deck, decrements waitgroup when done
	// wg.Add(1)
	// go playRound(&roundCopy)
	//go shuffleDeck(sorted, size)
	//}
	//Concurrently waits for all goroutines to finish to close the "sorted" channel
	go func() {
		wg.Wait()
		//close(sorted)
		close(roundsPlayed)
	}()
	//Loops through sorted channel incrementing true cases
	//Finishes loop when it receives "close" on the sorted channel
	for range roundsPlayed {
		//fmt.Printf("values are %v\n", v)
		combinations++
	}
	//Probability is number of true cases over trials
	fmt.Println("finished", combinations)
	//fmt.Printf("roundCopy after has these hands :: %v%v\n", roundCopy.CardsPlayed, roundCopy.Players)
	//fmt.Printf("startingRound after has these hands :: %v%v\n", startingRound.CardsPlayed, startingRound.Players)
	return
}

func playRound(round *Round) {
	wg.Add(1)
	defer wg.Done()
	// lastPlayer := len(round.Players) - 1
	//currentPlayer := round.Players[round.CurrentPosition]
	//var tempRound Round
	var currentHand []Card
	for {
		if round.Done() {
			roundsPlayed <- round.CardsPlayed
			return
		}
		//cardToPlay := 0
		currentHand = round.Players[round.CurrentPosition].Hand
		if len(currentHand) > 1 {
			for i := 1; i < len(currentHand); i++ {
				tempRound := round.CopyOf()
				tempRound.playNextCard(i)

				go playRound(&tempRound)
			}
			// for i := range currentHand[1:]{
			// 	tempRound := round.CopyOf()
			// 	tempRound.playNextCard(i)
			// 	wg.Add(1)
			// 	go playRound(&tempRound)
			// }
		}
		round.playNextCard(0)
	}

}
