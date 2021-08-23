# Initial stage: download modules
FROM golang:1.16-alpine as builder

WORKDIR /auth

COPY ./ /auth

RUN go mod download

# Intermediate stage: Build the binary
FROM golang:1.16-alpine as runner

COPY --from=builder ./auth ./auth

RUN go get github.com/githubnemo/CompileDaemon

WORKDIR /auth
ENV config=config

EXPOSE 5000

ENTRYPOINT CompileDaemon --build="go build cmd/auth/main.go" --command=./main
