build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main ./cmd/gh-app-pr-hello/...
	zip main.zip main


deploy:
	aws lambda update-function-code --function-name githubApp --zip-file fileb://main.zip

all: build deploy
