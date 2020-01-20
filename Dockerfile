####
# Build the binary
####
FROM golang:alpine AS builder-go

# Copy all source files.
WORKDIR /go/src/github.com/jheidel/rf64-convert/
COPY . .
# Build the executable.
RUN go build

####
# Install into the minimal image.
####
FROM alpine
WORKDIR /bin
COPY --from=builder-go /go/src/github.com/jheidel/rf64-convert/rf64-convert /bin
