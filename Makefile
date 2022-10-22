info:
	@echo 'make app binary: bin (b), release (r), debug (d)'
	@echo 'test with data: test'

all: debug release

done:
	@echo '- - -'
	@echo 'DONE.'

clear:
	clear

bin: clear go-compile done
b: bin

debug: clear go-debug done
d: debug

release: clear go-release done
r: release

test:
	./app/cmd -f ./files/alef.wav


go-compile:
	go build -v -o ./app/cmd ./cmd

go-debug:
	go build -v -tags="debug" -gcflags=all="-N -l" -o ./app/debug ./cmd

go-release:
	go build -v -a -tags="release" -ldflags='-s -w' -o ./app/cmd ./cmd