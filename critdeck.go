package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/oleiade/reflections"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Deck struct {
	Cards []Card `json:"Cards"`
}

func (deck *Deck) Random(number int, path string) (output []string) {
	for i := 0; i < number; i++ {
		output = append(output, deck.GetRandom(path))
	}
	return
}

func (deck *Deck) GetRandom(path string) (output string) {
	fmt.Println(path)
	card_num := rand.Intn(len(deck.Cards))
	card := deck.Cards[card_num]
	output = fmt.Sprint(reflections.GetField(card, path))
	return
}

type DeckSet struct {
	HitDeck Deck `json:"Hit"`
	MissDeck Deck `json:"Miss"`
}

type Card struct {
	Slashing string `json:"Slashing"`
	Bludgeoning string `json:"Bludgeoning"`
	Piercing string `json:"Piercing"`
	Magic string `json:"Magic"`
}


var decks DeckSet
func main() {
	rand.Seed(time.Now().UnixNano())
	jsonFile, err := ioutil.ReadFile("dek.json")
	if err != nil {
		log.Fatal("Could not read file")
	}
	json.Unmarshal(jsonFile, &decks)
	r := mux.NewRouter()
	r.HandleFunc("/{deck}", DrawCards).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", r))
}

func DrawCards(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	text := r.Form.Get("text")
	textSplit := strings.Split(text, " ")
	draw, _ := strconv.Atoi(textSplit[1])
	cardType := textSplit[0]
	vars := mux.Vars(r)
	fmt.Println(vars)
	deck := Deck{}
	if vars["deck"] == "Hit" {
		deck = decks.HitDeck
	} else if vars["deck"] == "Miss" {
		deck = decks.MissDeck
	}
	cards := deck.Random(draw, cardType)
	w.WriteHeader(http.StatusOK)
	for _, value := range cards {
		fmt.Fprintf(w, "%v\n", value)
	}
}