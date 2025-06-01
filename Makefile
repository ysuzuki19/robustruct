test:
	go test ./...

generate:
	go generate ./...

update:
	go list -m -f '{{if and (ne .Main true) (not .Indirect)}}{{.Path}}{{end}}' all | xargs -n1 go get -u=patch
	go mod tidy
