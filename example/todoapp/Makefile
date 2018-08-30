all:
	go install \
		github.com/netbrain/splitscreen/cmd/ss-aggregate \
		github.com/netbrain/splitscreen/cmd/ss-listener && \
	go generate github.com/netbrain/splitscreen/example/todoapp/todo/...
	go run "${GOPATH}/src/github.com/netbrain/splitscreen/example/todoapp/main.go"
