FROM golang:1.18-alpine AS build

# Environment Variables
ENV CGO_ENABLED 0
ENV GOOS linux
ENV GOARCH amd64

WORKDIR $HOME/github.com/JulesMike/spoty

# Add git, curl and upx support
RUN apk add --update git curl upx

# Pull modules
COPY go.* ./
RUN go mod download

# Copy code into image
COPY . .

# Build application for deployment
RUN go build -tags=jsoniter -trimpath -a -ldflags "-s -w -extldflags '-static'" \
    -v -o /spoty .

# Compress binary
RUN upx --best --lzma /spoty

# Create minimal image with just the application
# gcr.io/distroless/static is perfect for Go app that do not depend on libc
FROM gcr.io/distroless/static
COPY --from=build /spoty /
CMD ["/spoty", "serve"]
