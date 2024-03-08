FROM golang:1.22.1-alpine
WORKDIR /home/app
COPY . .
RUN GOOS=linux GOARCH=amd64 go build -o main ./main.go