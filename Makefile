all: install
	./bin/lear
install: main.go internal/text/lear.txt
	rm -f internal/index/*.idx
	go generate ./...
	go install ./
internal/text/lear.txt:
	curl https://gutenberg.org/cache/epub/1532/pg1532.txt > ./internal/text/lear.txt
clean:
	rm bin/lear
