run: t css sqlc
	@go build -o ./bin/out && ./bin/out

sqlc:
	@sqlc generate

gooseup:
	@cd models/sql/schema && goose postgres postgres://postgres:amin235711@amin-laptop.local:5432/daq up

goosedown:
	@cd models/sql/schema && goose postgres postgres://postgres:amin235711@amin-laptop.local:5432/daq down

css:
	@npx tailwindcss -i public/input.css -o public/style.css --minify

t:
	templ generate