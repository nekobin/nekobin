FROM golang:1.14-alpine

WORKDIR /app

# Copy go mod and sum files and download dependencies
COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

# Build the Go app to /go/bin/nekobin
RUN go build -o /go/bin/nekobin

EXPOSE 5555

CMD ["/go/bin/nekobin"]
