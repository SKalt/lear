package main

import (
	_ "embed"
	"encoding/binary"
	"fmt"
	"slices"
	"strings"
	"time"

	"gitub.com/skalt/lear/internal/text"
)

//go:generate go run ./internal/index/index.go

//go:embed internal/index/quotes.idx
var quoteIndex []byte

//go:embed internal/index/acts.idx
var actIndex []byte

//go:embed internal/index/scenes.idx
var sceneIndex []byte

var acts, scenes []int

func init() {
	scenes = make([]int, 0, len(sceneIndex)/4)
	for chunk := range slices.Chunk(sceneIndex, 4) {
		scenes = append(scenes, int(binary.BigEndian.Uint32(chunk)))
	}
	acts = make([]int, 0, len(actIndex)/4)
	for chunk := range slices.Chunk(actIndex, 4) {
		acts = append(acts, int(binary.BigEndian.Uint32(chunk)))
	}
}

var numerals = []string{
	"I", "II", "III", "IV", "V", "VI", "VII",
}

func romanize(i int) string {
	return numerals[i]
}

func fromIndex(idx []byte, i int) int {
	return int(binary.BigEndian.Uint32(idx[i*4 : i*4+4]))
}

func getQuote(tty bool) (result string) {
	output := strings.Builder{}
	nQuotes := len(quoteIndex)/4 - 2 // lead-up, trailing space
	quoteNumber := int(time.Now().UnixMilli() % int64(nQuotes))
	act, _ := slices.BinarySearch(acts, quoteNumber)
	scene, _ := slices.BinarySearch(scenes, quoteNumber)
	start := fromIndex(quoteIndex, quoteNumber)
	end := fromIndex(quoteIndex, quoteNumber+1)
	output.WriteString(fmt.Sprintf("KING LEAR Act %s, Scene %s\n", romanize(act), romanize(scene)))
	output.WriteString(text.Lear[start:end])
	output.WriteRune('\n')
	return output.String()
}

func main() {
	fmt.Println(getQuote(true))
}
