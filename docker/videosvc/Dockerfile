FROM golang:1.17-alpine AS builder

# Move to working directory (/videosvc).
WORKDIR /videosvc

# Copy the code into the container and download dependency.
COPY . .
RUN go mod download
# Set necessary environmet variables needed for our image and build the service.
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags="-s -w" -o videosvcbin cmd/videosvc/main.go


FROM jrottenberg/ffmpeg:3-scratch

# Copy binary file from /videosvc to root folder of scratch container.
COPY --from=builder ["/videosvc/videosvcbin", "/"]


# Command to run when starting the container.
ENTRYPOINT ["/videosvcbin"]
