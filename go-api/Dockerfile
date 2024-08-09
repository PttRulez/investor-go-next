FROM golang:1.22.0-alpine3.19

WORKDIR /app

COPY . .
COPY ./go.mod .
COPY ./go.sum .

# Download and install the dependencies:
RUN apk update && apk add git
RUN go get -d -v ./...

# Build the go app
RUN mv ./.env.docker ./.env
RUN go build -o controller ./cmd/app
EXPOSE 3001

CMD ["./api"]