# TODO:
# - secure Dockerfile

FROM golang:alpine AS builder

LABEL maintainer="Fabien Bellanger <valentil@gmail.com>"

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
# RUN go build -ldflags "-s -w" -a -installsuffix cgo -o fiber-boilerplate cmd/main.go
RUN go build -ldflags "-s -w" -a -installsuffix cgo -o fiber-boilerplate cmd/main.go

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN cp /build/fiber-boilerplate /build/favicon.png .
RUN cp -R /build/assets /build/templates .
RUN cp /build/.env.docker ./.env

# -----------------------------------------------------------------------------

FROM alpine:latest

LABEL maintainer="Fabien Bellanger <valentil@gmail.com>"

RUN apk update && apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /dist/.env .
COPY --from=builder /dist/favicon.png .
COPY --from=builder /dist/assets assets
COPY --from=builder /dist/templates templates
COPY --from=builder /dist/fiber-boilerplate .

RUN echo $(ls -lah)

# RUN chmod +x ./fiber-boilerplate

EXPOSE 3002
ENTRYPOINT ["./fiber-boilerplate"]
CMD ["run"]
