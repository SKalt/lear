all: install
	./bin/lear
install: main.go lear.txt
	go install ./
lear.txt:
	curl https://gutenberg.org/cache/epub/1532/pg1532.txt > ./lear.txt
clean:
	rm bin/lear
