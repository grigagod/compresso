FROM golang:1.16-alpine AS builder

# Move to working directory (/auth).
WORKDIR /auth

# Copy the code into the container and download dependency.
COPY . .
RUN go mod download
# Set necessary environmet variables needed for our image and build the API server.
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags="-s -w" -o authbin cmd/auth/main.go

FROM scratch

# Copy binary and config files from /auth to root folder of scratch container.
COPY --from=builder ["/auth/authbin", "/"]
COPY --from=builder ["/auth/docs/auth/", "/docs/auth/"]
COPY --from=builder ["/auth/config/auth/", "/config/auth/"]

ENV config-auth=docker

# Export necessary port.
EXPOSE 5000

# Command to run when starting the container.
ENTRYPOINT ["/authbin"]
