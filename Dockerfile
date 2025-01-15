# Build stage for UI
FROM node:18-alpine AS ui-builder
WORKDIR /stash
COPY . .
RUN apk add --no-cache git
WORKDIR /stash/ui/v2.5
RUN yarn install
RUN yarn gqlgen
RUN yarn build

# Build stage for Go
FROM golang:1.22-alpine AS builder
WORKDIR /stash
COPY . .
COPY --from=ui-builder /stash/ui/v2.5/build ui/v2.5/build
RUN apk add --no-cache build-base pkgconfig vips-dev
RUN go build -o stash cmd/stash/main.go

# Runtime stage
FROM alpine:latest
RUN apk add --no-cache \
    ffmpeg \
    vips \
    ca-certificates

WORKDIR /
COPY --from=builder /stash/stash /usr/local/bin/stash
COPY --from=builder /stash/ui/v2.5/build /usr/local/share/stash/ui

# Create necessary directories for VHS clips
RUN mkdir -p /generated/vhs_clips

EXPOSE 9999

ENTRYPOINT ["stash"]