package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Flashcard struct {
	Term          string `json:"term"`
	Definition    string `json:"definition"`
	MistakesCount int    `json:"mistakes_count"`
}

func addFlashcard(cards *[]Flashcard, terms map[string]bool, cardsByDefs map[string]Flashcard) {
	flashcard := Flashcard{}
	lg.Println("The card:")
	for {
		flashcard.Term = lg.ReadString()
		if terms[flashcard.Term] {
			lg.Printf("The term \"%s\" already exists. Try again:\n", flashcard.Term)
		} else {
			terms[flashcard.Term] = true
			break
		}
	}
	lg.Println("The definition of the card:")
	for {
		flashcard.Definition = lg.ReadString()
		if _, ok := cardsByDefs[flashcard.Definition]; ok {
			lg.Printf("The definition \"%s\" already exists. Try again:\n", flashcard.Definition)
		} else {
			cardsByDefs[flashcard.Definition] = flashcard
			break
		}
	}
	lg.Printf("The pair (\"%s\":\"%s\") has been added.\n", flashcard.Term, flashcard.Definition)
	*cards = append(*cards, flashcard)
}

func removeFlashcardByIdx(s []Flashcard, i int) []Flashcard {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func removeFlashcard(cards *[]Flashcard, terms map[string]bool, cardsByDefs map[string]Flashcard) {
	lg.Println("Which card?")
	term := lg.ReadString()
	if terms[term] {
		delete(terms, term)
		for k, v := range cardsByDefs {
			if v.Term == term {
				delete(cardsByDefs, k)
				break
			}
		}
		for i, v := range *cards {
			if v.Term == term {
				*cards = removeFlashcardByIdx(*cards, i)
				break
			}
		}
		lg.Println("The card has been removed.")
	} else {
		lg.Printf("Can't remove \"%s\": there is no such card.\n", term)
	}
}

func askCards(cards []Flashcard, cardsByDefs map[string]Flashcard) {
	if len(cards) == 0 {
		lg.Println("No cards available, please add or import first.")
		return
	}
	lg.Println("How many times to ask?")
	n, _ := strconv.Atoi(lg.ReadString())

	for i := 0; i < n; i++ {
		checkFlashcard(&cards[i%len(cards)], cardsByDefs)
	}
}

func checkFlashcard(flashcard *Flashcard, cardsByDefs map[string]Flashcard) {

	lg.Printf("Print the definition of \"%s\":\n", flashcard.Term)
	answer := lg.ReadString()
	if answer == flashcard.Definition {
		lg.Println("Correct!")
	} else {
		flashcard.MistakesCount++
		flashcardForDef, ok := cardsByDefs[answer]
		if ok {
			lg.Printf("Wrong. The right answer is \"%s\", ", flashcard.Definition)
			lg.Printf("but your definition is correct for \"%s\".\n", flashcardForDef.Term)
		} else {
			lg.Printf("Wrong. The right answer is \"%s\".\n", flashcard.Definition)
		}
	}
}

func exportFlashcardsToFilename(filename string, cards []Flashcard) {
	cardsJson, err := json.Marshal(cards)
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(filename, []byte(cardsJson), 0664)
	if err != nil {
		log.Fatal(err)
	}
	lg.Printf("%d cards have been saved.\n", len(cards))
}

func exportFlashcards(cards []Flashcard) {
	lg.Println("File name:")
	filename := lg.ReadString()

	exportFlashcardsToFilename(filename, cards)
}

func importFlashcardsFromFile(filename string, cards *[]Flashcard, terms *map[string]bool, cardsByDefs *map[string]Flashcard) {
	cardsJson, err := os.ReadFile(filename)
	if err != nil {
		lg.Println("File not found.")
		return
	}

	var newCards []Flashcard
	err = json.Unmarshal(cardsJson, &newCards)
	if err != nil {
		log.Fatal(err)
	}

	lg.Printf("%d cards have been loaded.\n", len(newCards))

	*terms = make(map[string]bool)
	*cardsByDefs = make(map[string]Flashcard)
	for _, card := range newCards {
		(*terms)[card.Term] = true
		(*cardsByDefs)[card.Definition] = card
	}
	for _, card := range *cards {
		if !(*terms)[card.Term] {
			_, ok := (*cardsByDefs)[card.Definition]
			if !ok {
				newCards = append(newCards, card)
				(*terms)[card.Term] = true
				(*cardsByDefs)[card.Definition] = card
			}
		}
	}

	*cards = newCards
}

func importFlashcards(cards *[]Flashcard, terms *map[string]bool, cardsByDefs *map[string]Flashcard) {
	lg.Println("File name:")
	filename := lg.ReadString()

	importFlashcardsFromFile(filename, cards, terms, cardsByDefs)
}

func hardestCard(cards []Flashcard) {
	var hardest []Flashcard
	for _, card := range cards {
		if card.MistakesCount > 0 {
			if hardest == nil || card.MistakesCount > hardest[0].MistakesCount {
				hardest = []Flashcard{card}
			} else if card.MistakesCount == hardest[0].MistakesCount {
				hardest = append(hardest, card)
			}
		}
	}

	if hardest == nil {
		lg.Println("There are no cards with errors.")
	} else if len(hardest) == 1 {
		lg.Printf(
			"The hardest card is \"%s\". You have %d errors answering it.\n",
			hardest[0].Term,
			hardest[0].MistakesCount)
	} else {
		terms := make([]string, len(hardest))
		for i, card := range hardest {
			terms[i] = "\"" + card.Term + "\""
		}
		lg.Printf(" The hardest cards are %s\n", strings.Join(terms, ", "))
	}
}

func resetStats(cards []Flashcard) {
	for i, _ := range cards {
		cards[i].MistakesCount = 0
	}
	lg.Println("Card statistics have been reset.")
}

func saveLog() {
	lg.Println("File name:")
	filename := lg.ReadString()

	lg.Save(filename)

	lg.Println("The log has been saved.")
}

type Log struct {
	strings.Builder
	reader *bufio.Reader
}

func (log *Log) Println(v ...any) {
	message := fmt.Sprintln(v...)
	log.WriteString(message)
	fmt.Print(message)
}

func (log *Log) Printf(format string, v ...any) {
	message := fmt.Sprintf(format, v...)
	log.WriteString(message)
	fmt.Print(message)
}

func (log *Log) ReadString() string {
	s, _ := log.reader.ReadString('\n')
	log.WriteString(s)
	return strings.TrimSpace(s)
}

func (log *Log) Save(filename string) {
	os.WriteFile(filename, []byte(log.String()), 0664)
}

var lg = Log{strings.Builder{}, bufio.NewReader(os.Stdin)}

func main() {
	cards := make([]Flashcard, 0)
	terms := make(map[string]bool, 0)
	cardsByDefs := make(map[string]Flashcard, 0)

	importFilename := flag.String("import_from", "", "the filename to import from on initial start")
	exportFilename := flag.String("export_to", "", "the filename to export before the program is finished")

	flag.Parse()

	if *importFilename != "" {
		importFlashcardsFromFile(*importFilename, &cards, &terms, &cardsByDefs)
	}

	action := ""
	for action != "exit" {
		lg.Println("Input the action (add, remove, import, export, ask, exit, log, hardest card, reset stats):")
		action = lg.ReadString()

		switch action {
		case "add":
			addFlashcard(&cards, terms, cardsByDefs)
		case "remove":
			removeFlashcard(&cards, terms, cardsByDefs)
		case "import":
			importFlashcards(&cards, &terms, &cardsByDefs)
		case "export":
			exportFlashcards(cards)
		case "ask":
			askCards(cards, cardsByDefs)
		case "log":
			saveLog()
		case "hardest card":
			hardestCard(cards)
		case "reset stats":
			resetStats(cards)
		}
		lg.Println("")
	}

	if *exportFilename != "" {
		exportFlashcardsToFilename(*exportFilename, cards)
	}
	lg.Println("Bye bye!")
}
