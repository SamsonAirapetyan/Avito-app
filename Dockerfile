FROM golang:1.21-alpine

RUN go version
ENV GOPATH=/

COPY ./ ./

# build go app
RUN go mod download
RUN go build -o Avito-app ./cmd/app/main.go

CMD ["./Avito-app"]