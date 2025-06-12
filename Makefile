BINARYNAME?=rickrolld
PWD?=`pwd`

.PHONY: mod rickrolld

mod:
	go get -u github.com/nugget/go-telnet
	go mod tidy
	go mod verify
	git diff go.mod go.sum && git commit go.mod go.sum -m "make mod"

go-telnet-local:
	-git clone git@github.com:nugget/go-telnet.git

localdev: go-telnet-local
	go mod edit -replace=github.com/nugget/go-telnet="$(PWD)/go-telnet"
	go mod tidy

productiondev:
	go mod edit -dropreplace=github.com/nugget/go-telnet
	go mod tidy

rickrolld: 
	mkdir -p dist
	clearbuffer && go mod tidy && go build -o dist/$(BINARYNAME) ./rickrolld/main.go
	./dist/$(BINARYNAME) -v

