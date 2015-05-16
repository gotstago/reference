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
	north := shuffledCards[0:3]
	t.Logf("North Cards are %v, capacity is %d", north, cap(north))
	east := shuffledCards[9:12]
	t.Logf("East Cards are %v, capacity is %d", east, cap(east))
	south := shuffledCards[18:21]
	t.Logf("South Cards are %v, capacity is %d", south, cap(south))
	west := shuffledCards[27:30]
	t.Logf("West Cards are %v, capacity is %d", west, cap(west))
	//http://stackoverflow.com/questions/25025409/delete-element-in-a-slice
	west = append(west[:1], west[2:]...)
	t.Logf("West Cards are %v, capacity is %d", west, cap(west))

	//}
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
