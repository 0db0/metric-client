FROM golang:1.23-alpine

ARG USER_ID
ARG GROUP_ID

ENV GOFLAGS="-buildvcs=false"

WORKDIR /app

RUN apk add --no-cache git \
    && go install github.com/air-verse/air@v1.61.5 \
    && go install github.com/go-delve/delve/cmd/dlv@v1.24.0

COPY go.mod go.sum ./
RUN go mod download

COPY . .

CMD ["air", "-c", ".air.toml"]
