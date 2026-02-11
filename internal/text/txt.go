package text

import (
	_ "embed"
	"strings"
)

//go:embed lear.txt
var Lear string

const start = "SCENE: Britain"

const end = "*** END OF THE PROJECT GUTENBERG EBOOK KING LEAR ***"

func init() {
	// pre-process the text
	Lear = Lear[strings.Index(Lear, start)+len(start) : strings.Index(Lear, end)]
}
