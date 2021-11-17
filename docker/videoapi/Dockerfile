FROM golang:1.17-alpine AS builder

# Move to working directory (/videoapi).
WORKDIR /videoapi

RUN apk add -U --no-cache ca-certificates

# Copy the code into the container and download dependency.
COPY . .
RUN go mod download
# Set necessary environmet variables needed for our image and build the API server.
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags="-s -w" -o videoapibin cmd/videoapi/main.go


FROM scratch

# Copy binary and documentation files from /videoapi to root folder of scratch container.
COPY --from=builder ["/videoapi/videoapibin", "/"]
COPY --from=builder ["/videoapi/docs/videoapi/", "/docs/videoapi/"]

# Copy certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/


# Command to run when starting the container.
ENTRYPOINT ["/videoapibin"]
