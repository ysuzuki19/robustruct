test:
	go test ./...

generate:
	go generate ./...

update:
	@go list -m -f '{{if and (ne .Main true) (not .Indirect)}}{{.Path}}{{end}}' all | xargs -n1 go get -u=patch
	@go mod tidy

outdated: outdated-patch outdated-minor

outdated-patch::
	@echo "================================"
	@echo "- patch update -----------------"
	@echo "                                "
	@go list -m -u -json all | jq -r 'select(.Main != true and .Indirect != true and .Update) | select( (.Version | split(".")[:2]) == (.Update.Version | split(".")[:2])) | "\(.Path): \(.Version) → \(.Update.Version)"'
	@echo "--------------------------------"
	@echo "Run 'make update-patch' to update the dependencies."
	@echo "================================"

outdated-minor:
	@echo "================================"
	@echo "- minor update -----------------"
	@echo "--------------------------------"
	@go list -m -u -f '{{if and (not .Main) (not .Indirect) .Update}}{{.Path}}: {{.Version}} → {{.Update.Version}}{{end}}' all
	@echo "--------------------------------"
	@echo "================================"
