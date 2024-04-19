dev:
	air -c .air.toml
serve:
	templ generate && go run . serve
build:
	templ generate && CGO_ENABLED=0 go build -o ./main
kill:
	pkill 8090
types:
	pb-gen models
clean:
	rm -rf main temp