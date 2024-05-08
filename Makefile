dev:
	air -c .air.toml
serve:
	templ generate && go run . serve
build:
	templ generate && CGO_ENABLED=0 go build -o ./main
kill:
	pkill 8090
types:
	go run . types
styles:
	npx tailwindcss -i ./public/styles.css -o ./public/tailwind.css
clean:
	rm -rf main temp