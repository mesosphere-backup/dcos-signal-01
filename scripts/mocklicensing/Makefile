PIDFILE = .serve-license.pid
BINARY  = mocklicensing

.PHONY: build clean start stop request $(PIDFILE)

build: mock_license_server.go
	@go build -o $(BINARY) ./...

clean:
	@rm $(BINARY)

request:
	@curl -s --unix-socket /tmp/dcos-licensing.socket http: | jq .

start: $(PIDFILE)
	@echo "Running mocklicensing..."

$(PIDFILE):
	@./$(BINARY) & echo $$! > $@

stop: $(PIDFILE)
	@kill `cat $<` && rm $< && echo "Stopping mocklicensing..."

