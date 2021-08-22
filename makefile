bin/arb: .force
	go build -o bin/arb .
bin/arb-linux: .force
	env GOOS=linux GOARCH=386 go build -v -o bin/arb_linux
update-eos-go:
	go get -u github.com/panyanyany/eos-go@develop
.force:
