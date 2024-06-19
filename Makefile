build:
	go build -o ./cmd/bin/shortener ./cmd/shortener/main.go

run:
	CONFIG="./configs/local_config.yaml" PRIVATE_KEY="super secret key for tokens" ./cmd/bin/shortener

docker:
	docker build -t shortener-docker