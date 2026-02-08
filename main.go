package main

import (
	_ "embed"
	"fmt"
	"regexp"
	"strings"
	"time"
)

var quotePattern = regexp.MustCompile(`[A-Z]+\.`)

const (
	lines = 1051
	// acts   = 10
	// scenes = 26
)
const start = "*** START OF THE PROJECT GUTENBERG EBOOK KING LEAR ***"
const end = "*** END OF THE PROJECT GUTENBERG EBOOK KING LEAR ***"

//go:embed lear.txt
var lear string
var quoteRanges [][]int

func init() {
	// pre-process the text
	lear = lear[strings.Index(lear, start)+len(start) : strings.Index(lear, end)]
	quoteRanges = quotePattern.FindAllStringIndex(lear, lines)
	if len(quoteRanges) != lines {
		panic(fmt.Errorf("nQuotes %d != %d lines", len(quoteRanges), lines))
	}
}

func main() {
	quoteNumber := int(time.Now().UnixMilli() % lines)
	q := quoteRanges[quoteNumber] // random line
	r := []int{len(lear)}
	// q1 := len(lear)
	if quoteNumber < len(quoteRanges) {
		r = quoteRanges[quoteNumber+1]
	}
	fmt.Println(lear[q[0]:r[0]])
}
