FROM golang:1.22.0-alpine3.19

WORKDIR /app

COPY . .
COPY ./go.mod .
COPY ./go.sum .

# Download and install the dependencies:
RUN go get -d -v ./...

# Build the go app
RUN mv ./.env.docker ./.env
RUN go build -o api ./cmd/app
EXPOSE 3001

CMD ["./api"]