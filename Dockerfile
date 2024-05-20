FROM golang:1.22 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /post-comment-system ./cmd/server

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=builder /post-comment-system /post-comment-system

EXPOSE 8080

CMD ["/post-comment-system"]