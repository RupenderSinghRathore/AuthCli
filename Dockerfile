FROM golang:1.26-alpine AS builder

WORKDIR /src

RUN apk add --no-cache build-base git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=1

RUN go build -o /out/authcli ./cmd/app
RUN go install github.com/amacneil/dbmate/v2@latest

FROM alpine:3.22

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /out/authcli /app/authcli
COPY --from=builder /go/bin/dbmate /usr/local/bin/dbmate
COPY db/migrations ./db/migrations
COPY docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh

RUN chmod +x /usr/local/bin/docker-entrypoint.sh

ENV DATABASE_PATH=/app/db/myapp.db
ENV DATABASE_URL=sqlite:///app/db/myapp.db

ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]
CMD ["/app/authcli"]
