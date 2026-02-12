package main

import (
	_ "embed"
	"encoding/binary"
	"fmt"
	"os"
	"slices"
	"strconv"
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
	scenes = make([]int, 0, len(sceneIndex)/2)
	for chunk := range slices.Chunk(sceneIndex, 2) {
		scenes = append(scenes, int(binary.BigEndian.Uint16(chunk)))
	}
	acts = make([]int, 0, len(actIndex))
	for _, a := range actIndex {
		acts = append(acts, int(uint8(a)))
	}
}

var numerals = []string{
	"", "I", "II", "III", "IV", "V", "VI", "VII",
}

func romanize(i int) string {
	return numerals[i]
}

func fromIndex(idx []byte, i int) int {
	return int(binary.BigEndian.Uint32(idx[i*4 : i*4+4]))
}

func getAct(sceneNumber int) (act int) {
	return binarySearchRoundDown(acts, sceneNumber)
}
func binarySearchRoundDown(arr []int, x int) (result int) {
	result, _ = slices.BinarySearch(arr, x)
	if arr[result] > x {
		result -= 1
	}
	return
}
func getScene(quoteNumber int) (scene int) {
	return binarySearchRoundDown(scenes, quoteNumber)
}
func getLoc(quoteNumber int) (act, scene int) {
	scene = getScene(quoteNumber)
	act = getAct(scene)
	// hack: adjust scene
	if act == 0 {
		scene += 1
	}
	return
}

var nQuotes = len(quoteIndex)/4 - 2 // lead-up, trailing space
func getQuote(quoteNumber int) string {
	start := fromIndex(quoteIndex, quoteNumber)
	end := fromIndex(quoteIndex, quoteNumber+1)
	return text.Lear[start:end]
}
func must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}
func main() {
	var quoteNumber int
	if len(os.Args) > 1 {
		quoteNumber = must(strconv.Atoi(os.Args[1]))
	} else {
		quoteNumber = int(time.Now().UnixMilli() % int64(nQuotes))
	}
	fmt.Println(quoteNumber)
	quote := getQuote(quoteNumber)
	output := strings.Builder{}
	act, scene := getLoc(quoteNumber)

	fmt.Fprintf(&output, "KING LEAR Act %s",
		romanize(act+1),
	)
	if scene > 0 {
		fmt.Fprintf(&output, ", Scene %s", romanize(scene-acts[act]))
	}
	output.WriteString("\n")
	fmt.Fprint(&output, quote)
	fmt.Println(output.String())
}
