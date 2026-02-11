#!/bin/sh
awk '
  /^[A-Z]+\./ { lines++ }
  /^ACT /    { acts++  }
  /^SCENE /  { scenes++}
  END  {
    print "lines = " lines
    print "acts = "  acts
    print "scenes  = " scenes
  }
' ./internal/text/lear.txt

# printf "number of quotes: "
# grep -e '^[A-Z]+\.'  | wc -l

# printf "number of acts: "
# grep -e '^ACT ' lear.txt | wc -l

# printf "number of schenes: "
# grep -e '^SCENE '
