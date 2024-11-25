FROM golang:1.23

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY / /app

RUN CGO_ENABLED=1 GOOS=linux go build -o /bsky-feed-server

EXPOSE 8081

CMD [ "/bsky-feed-server" ]