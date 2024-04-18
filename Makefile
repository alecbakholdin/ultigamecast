dev:
	wgo -file=.go -file=.templ -xfile=_templ.go templ generate :: go run . serve
serve:
	templ generate && go run . serve
build:
	templ generate && CGO_ENABLED=0 go build -o ./main
types:
	pb-gen models
clean:
	rm -rf main temp