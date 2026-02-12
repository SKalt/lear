// generate an array of lengths of quotes
package main

import (
	"encoding/binary"
	"os"
	"regexp"
	"strings"

	"github.com/skalt/pathlib.go"
	"gitub.com/skalt/lear/internal/text"
)

var quoteP = regexp.MustCompile(`^[A-Z]+\.`)

type Q int // quote-index
type S int // scene index
var Lear string = text.Lear

func must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}

// var romanNumeralPattern = regexp.MustCompile(`^[IVX]+.$`)
func main() {
	// 0th quote is the lead-up to the first bit
	quoteRanges := make([]int, 0, 1200)
	scenes := make([]Q, 1, 31)
	acts := make([]S, 0, 4)
	var i int      // in ascii characters
	var sceneLen Q // in quotes
	var actLen S   // in scenes
	nextQuote := func() {
		quoteRanges = append(quoteRanges, i)
		sceneLen += Q(1)
	}
	nextScene := func() {
		nextQuote() // terminate any ongoing quote
		scenes = append(scenes, sceneLen)
		actLen += S(1)
	}
	nextAct := func() {
		nextScene()
		acts = append(acts, actLen)
	}
	lines := strings.SplitSeq(Lear, "\n")
	// 0th quote: lead-up to the first bit
	for line := range lines {
		if quoteP.MatchString(line) {
			nextAct()
			break
		} else {
			i += len(line) + 1
		}
	}
	for line := range lines {
		if quoteP.MatchString(line) {
			nextQuote()
		} else if strings.HasPrefix(line, "SCENE ") {
			nextScene()
		} else if strings.HasPrefix(line, "ACT ") {
			nextAct()
		}
		i += len(line) + 1
	}
	nextAct() // finish
	// quoteRanges = append(quoteRanges, len(Lear))
	// scenes = append(scenes, sceneLen)
	// acts = append(acts, actLen)
	// cwd := must(pathlib.Cwd())
	cwd := must(pathlib.Dir("~/programming/lear").ExpandUser())
	thisDir := cwd.Join("internal/index").AsDir()
	{
		quoteIndex := make([]byte, 0, 4*len(quoteRanges))
		for _, q := range quoteRanges {
			quoteIndex = binary.BigEndian.AppendUint32(quoteIndex, uint32(q))
		}
		target := must(thisDir.Join("quotes.idx").AsFile().Open(os.O_CREATE|os.O_RDWR, 0o0666))
		defer target.Close()
		_ = must(target.Write(quoteIndex))
	}
	{
		sceneIndex := make([]byte, 0, 2*len(scenes))
		for _, s := range scenes {
			sceneIndex = binary.BigEndian.AppendUint16(sceneIndex, uint16(s))
		}
		target := must(thisDir.Join("scenes.idx").AsFile().Open(os.O_CREATE|os.O_RDWR, 0o666))
		defer target.Close()
		_ = must(target.Write(sceneIndex))
	}
	{
		actIndex := make([]byte, 0, len(acts))
		for _, a := range acts {
			actIndex = append(actIndex, uint8(a))
		}
		target := must(thisDir.Join("acts.idx").AsFile().Open(os.O_CREATE|os.O_RDWR, 0o0666))
		defer target.Close()
		_ = must(target.Write(actIndex))
	}
}
