package utilities

import (
	"fmt"
	//"github.com/gotstago/card"
	"github.com/gotstago/deck"
	// "net/http/httptest"
	"testing"
)

type PlayingCard struct {
	Rank string
	Suit string
}

func (p *PlayingCard) String() string {
	return fmt.Sprintf("%s of %s", p.Rank, p.Suit)
}

func TestSpecificDeck(t *testing.T) {
	//d := deck.NewDeck(false)
	d := deck.NewSpecificDeck(true, deck.FACES, []deck.Suit{deck.SPADE})
	//d.cards = append(d.cards, deck.Card{ACE, HEART}, deck.Card{KING, HEART})
	//result := fmt.Sprintf("%s", d)
	t.Logf("Number of Cards is %d", d.NumberOfCards())
	//assert.Equal(t, "A♥\nK♥\n", result, "These should be equal")
}

func TestPermutations(t *testing.T) {
	faces := []deck.Face{deck.ACE, deck.KING, deck.QUEEN, deck.JACK, deck.TEN, deck.NINE, deck.EIGHT, deck.SEVEN, deck.SIX}
	suits := deck.SUITS
	cards := make([]deck.Card, len(suits)*len(faces))
	for sindex, s := range suits {
		for findex, f := range faces {
			index := (sindex * len(faces)) + findex
			cards[index] = deck.Card{f, s}
		}
	}
	t.Logf("Number of Cards is %d", len(cards))
	d := deck.Deck{cards}
	//if shuffled {
	d.Shuffle()
	shuffledCards := d.Cards
	t.Logf("Cards are %v", shuffledCards)
	north := shuffledCards[0:2]
	t.Logf("North Cards are %v, capacity is %d", north, cap(north))
	east := shuffledCards[9:11]
	t.Logf("East Cards are %v, capacity is %d", east, cap(east))
	south := shuffledCards[18:20]
	t.Logf("South Cards are %v, capacity is %d", south, cap(south))
	west := shuffledCards[27:29]
	t.Logf("West Cards are %v, capacity is %d", west, cap(west))
	//http://stackoverflow.com/questions/25025409/delete-element-in-a-slice
	//west = append(west[:1], west[2:]...)
	//west = remove(0, west)
	//t.Logf("West Cards are %v, capacity is %d", west, cap(west))

	allCardsInRound := [][]deck.Card{north, east, south, west}
	t.Logf("all cards :: %v", allCardsInRound)

	for combination := range GenerateAllCombinations(allCardsInRound) {
		t.Log(combination) // This is instead of process(combination)
	}

	t.Log("Done!")
	/*for _, h := range allCardsInRound {
		for i, cell := range h {
			t.Logf("card is %v at position %d", cell, i)
		}
		t.Log("looping ...")
	}*/

	// for i, h := range allCardsInRound {
	// 	/*for i, cell := range h {
	// 		t.Logf("card is %v at position %d", cell, i)
	// 	}*/
	// 	t.Logf("Length before removal of first is %d, capacity is %d", len(h), cap(h))
	// 	allCardsInRound[i] = remove(0, allCardsInRound[i])
	// 	//t.Logf("card is %v at position 0", h[0])
	// 	t.Logf("Length after removal of first is %d, capacity is %d", len(h), cap(h))
	// 	t.Log("looping ...")
	// }
	// t.Logf("all cards :: %v", allCardsInRound)
	//}
}

func playRound(nextCardToPlay int, allCards [][]deck.Card, currentResult []deck.Card) {

}

func GenerateAllCombinations(allCards [][]deck.Card) <-chan []deck.Card {
	c := make(chan []deck.Card)

	// Starting a separate goroutine that will create all the combinations,
	// feeding them to the channel c
	go func(c chan []deck.Card) {
		defer close(c) // Once the iteration function is finished, we close the channel
		playedHand := make([]deck.Card, 0)
		NextPlay(c, 0, allCards, playedHand) // We start by feeding it 1st slice of cards
	}(c)

	return c // Return the channel to the calling function
}

// AddLetter adds a letter to the combination to create a new combination.
// This new combination is passed on to the channel before we call AddLetter once again
// to add yet another letter to the new combination in case length allows it
func NextPlay(c chan []deck.Card, index int, hands [][]deck.Card, played []deck.Card) {
	// Check if we reached the length limit
	// If so, we just return without adding anything
	if len(hands[index]) == 1 { ///*|| len(played) == len(hands)*len(hands[0]*/
		c <- played
		return
	}

	//var newCombo string
	for i, card := range hands[index] {
		copyOfPlayed := append([]deck.Card(nil), played...)
		copyOfPlayed = append(copyOfPlayed, card)
		copyOfHands := append([][]deck.Card(nil), hands...)
		copyOfHands[index] = remove(i, copyOfHands[index])
		NextPlay(c, (index+1)%4, copyOfHands, copyOfPlayed)
		// newCombo = combo + string(ch)
		// if len(newCombo) == 4 {
		// 	c <- newCombo
		// }
		// AddLetter(c, newCombo, alphabet, length-1)
	}
}

//////////////////////////////
//from http://stackoverflow.com/questions/19249588/go-programming-generating-combinations
func GenerateCombinations(alphabet string, length int) <-chan string {
	c := make(chan string)

	// Starting a separate goroutine that will create all the combinations,
	// feeding them to the channel c
	go func(c chan string) {
		defer close(c) // Once the iteration function is finished, we close the channel

		AddLetter(c, "", alphabet, length) // We start by feeding it an empty string
	}(c)

	return c // Return the channel to the calling function
}

// AddLetter adds a letter to the combination to create a new combination.
// This new combination is passed on to the channel before we call AddLetter once again
// to add yet another letter to the new combination in case length allows it
func AddLetter(c chan string, combo string, alphabet string, length int) {
	// Check if we reached the length limit
	// If so, we just return without adding anything
	if length <= 0 {
		return
	}

	var newCombo string
	for _, ch := range alphabet {
		newCombo = combo + string(ch)
		if len(newCombo) == 4 {
			c <- newCombo
		}
		AddLetter(c, newCombo, alphabet, length-1)
	}
}

//from http://stackoverflow.com/questions/19249588/go-programming-generating-combinations
////////////////////////////
func remove(element int, source []deck.Card) []deck.Card {
	source = append(source[:element], source[element+1:]...)
	return source
}

func TestTarabishSpecificDeck(t *testing.T) {
	//d := deck.NewDeck(false)
	d := deck.NewSpecificDeck(true,
		[]deck.Face{deck.ACE, deck.KING, deck.QUEEN, deck.JACK, deck.TEN, deck.NINE, deck.EIGHT, deck.SEVEN, deck.SIX},
		deck.SUITS)
	//d.cards = append(d.cards, deck.Card{ACE, HEART}, deck.Card{KING, HEART})
	//result := fmt.Sprintf("%s", d)
	t.Logf("Number in Deck is %d", d.NumberOfCards())
	//assert.Equal(t, "A♥\nK♥\n", result, "These should be equal")
}

func TestCards(t *testing.T) {
	t.Log("starting TestCards...")
	c1 := PlayingCard{"4", "h"}
	if c1.String() != "4 of h" {
		t.Error("Error printing card.")
	}
	//fmt.Println("begin test cards...")
	//t.Error("logging...")
	/*server := httptest.NewServer(new(HelloHandler))
	defer server.Close()

	// Pretend this is some sort of Go client...
	url := fmt.Sprintf("%s?say=Nothing", server.URL)
	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		t.Errorf("Error performing request.")
	}

	if resp.StatusCode != 404 {
		t.Errorf("Did not get a 404.")
	}*/
}

/*func TestTextHandler(t *testing.T) {
	handler := new(TextHandler)
	expectedBody := `
John Smith is 22 years old.

Alice Smith is 25 years old.

Bob Baker is 24 years old.
`

	recorder := httptest.NewRecorder()

	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost/hello?say=%s", expectedBody), nil)
	if err != nil {
		t.Errorf("Failed to create request.")
	}

	handler.ServeHTTP(recorder, req)

	switch recorder.Body.String() {
	case expectedBody:
		// body is equal so no need to do anything
	default:
		t.Errorf("Body (%s) did not match expectation (%s).",
			recorder.Body.String(),
			expectedBody)
	}
}

func TestEchosContent(t *testing.T) {
	handler := new(HelloHandler)
	expectedBody := "hellooo!"

	recorder := httptest.NewRecorder()

	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost/hello?say=%s", expectedBody), nil)
	if err != nil {
		t.Errorf("Failed to create request.")
	}

	handler.ServeHTTP(recorder, req)

	switch recorder.Body.String() {
	case expectedBody:
		// body is equal so no need to do anything
	default:
		t.Errorf("Body (%s) did not match expectation (%s).",
			recorder.Body.String(),
			expectedBody)
	}
}

func TestReturns404IfYouSayNothing(t *testing.T) {
	handler := new(HelloHandler)

	recorder := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "http://example.com/echo?say=Nothing", nil)
	if err != nil {
		t.Errorf("Failed to create request.")
	}

	handler.ServeHTTP(recorder, req)

	if recorder.Code != 404 {
		t.Errorf("Did not get a 404.")
	}
}

func TestClient(t *testing.T) {
	server := httptest.NewServer(new(HelloHandler))
	defer server.Close()

	// Pretend this is some sort of Go client...
	url := fmt.Sprintf("%s?say=Nothing", server.URL)
	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		t.Errorf("Error performing request.")
	}

	if resp.StatusCode != 404 {
		t.Errorf("Did not get a 404.")
	}
}*/
