FROM golang:1.23-alpine AS builder

WORKDIR /app

ENV GO111MODULE=on


COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o sql-mcp-server .


FROM alpine:latest

COPY --from=builder /app/sql-mcp-server /usr/local/bin/sql-mcp-server

EXPOSE 8088

CMD ["sql-mcp-server"]

