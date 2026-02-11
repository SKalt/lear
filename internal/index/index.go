// generate an array of lengths of quotes
package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"

	"github.com/skalt/pathlib.go"
	"gitub.com/skalt/lear/internal/text"
)

var quoteP = regexp.MustCompile(`^[A-Z]+\.`)

type Q int // quote-index
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
	scenes := make([]Q, 0, 30) // in
	acts := make([]Q, 0, 3)
	var i int      // in ascii characters
	var sceneLen Q // in quotes
	var actLen Q   // in scenes
	lines := strings.SplitSeq(Lear, "\n")
	// 0th quote: lead-up to the first bit
	for line := range lines {
		if quoteP.MatchString(line) {
			quoteRanges = append(quoteRanges, i)
			break
		} else {
			i += len(line) + 1
		}
	}
	for line := range lines {
		if quoteP.MatchString(line) {
			quoteRanges = append(quoteRanges, i)
			actLen += Q(1)
			sceneLen += Q(1)
		} else if strings.HasPrefix(line, "ACT ") {
			quoteRanges = append(quoteRanges, i)
			acts = append(acts, actLen)
		} else if strings.HasPrefix(line, "SCENE ") {
			quoteRanges = append(quoteRanges, i)
			scenes = append(scenes, sceneLen)
		}
		i += len(line) + 1
	}
	quoteRanges = append(quoteRanges, len(Lear))
	cwd := must(pathlib.Cwd())
	thisDir := cwd.Join("internal/index").AsDir()
	{
		quoteIndex := make([]byte, 0, 2*len(quoteRanges))
		for _, q := range quoteRanges {
			quoteIndex = binary.BigEndian.AppendUint32(quoteIndex, uint32(q))
		}
		target := must(thisDir.Join("quotes.idx").AsFile().Open(os.O_CREATE|os.O_RDWR, 0o0666))
		defer target.Close()
		_ = must(target.Write(quoteIndex))
	}
	{
		actIndex := make([]byte, 0, 2*len(acts))
		for _, a := range acts {
			actIndex = binary.BigEndian.AppendUint32(actIndex, uint32(a))
		}
		target := must(thisDir.Join("acts.idx").AsFile().Open(os.O_CREATE|os.O_RDWR, 0o0666))
		defer target.Close()
		_ = must(target.Write(actIndex))
	}

	{
		sceneIndex := make([]byte, 0, 2*len(scenes))
		for _, s := range scenes {
			sceneIndex = binary.BigEndian.AppendUint32(sceneIndex, uint32(s))
		}
		target := must(thisDir.Join("scenes.idx").AsFile().Open(os.O_CREATE|os.O_RDWR, 0o666))
		defer target.Close()
		_ = must(target.Write(sceneIndex))
	}
	fmt.Println(quoteRanges)
	fmt.Println(scenes)
	fmt.Println(acts)
	for i, q := range quoteRanges[1:] {
		if d := (q - quoteRanges[i]); d > 1000 {
			fmt.Println(i, d)

			act, _ := slices.BinarySearch(acts, Q(i))
			scene, _ := slices.BinarySearch(scenes, Q(i))
			fmt.Println("-------------------------------------------------------------")
			fmt.Println("Act", act+1, " scene ", scene+1)
			fmt.Println(Lear[quoteRanges[i]:q])
		}
	}
	// for i, q := range quoteRanges[1:] {
	// 	act, ok := slices.BinarySearch(acts, Q(i))
	// 	scene, ok := slices.BinarySearch(scenes, Q(i))
	// 	fmt.Println("-------------------------------------------------------------")
	// 	fmt.Println("Act", act, ok, " scene ", scene, ok)
	// 	fmt.Println(Lear[quoteRanges[i]:q])
	// }
	// todos:
	//  - romanize acts, scenes
	//  - scenes as offset from prev act
	//   - serialize to be bytes, embed in main
}
