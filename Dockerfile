FROM golang:1.11.2

COPY . /app
WORKDIR /app/bin
RUN go build ../cmd/wikipedia-extract/wikipedia-extract.go
RUN go build ../cmd/wikipedia-insert/wikipedia-insert.go
CMD ./import
