# Stage 1: build the Go binary
FROM golang:1.24-alpine AS builder

WORKDIR /src

# copy go.mod first so Docker can cache this step separately from source changes
COPY go.mod ./
RUN go mod download

COPY . .

# CGO_ENABLED=0 -> no C deps; -ldflags="-s -w" -> smaller binary
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o server .

# Stage 2: minimal final image, only what's needed to run
FROM alpine:3.20

# bash so `docker exec -it <container> /bin/bash` works (the audit uses this exact command)
RUN apk add --no-cache bash

LABEL maintainer="mozuhair" \
      description="ascii-art-web: Go web server that renders text as ASCII banners"

WORKDIR /app

# banner/*.txt is not copied — it's already embedded in the binary via go:embed
COPY --from=builder /src/server .
COPY --from=builder /src/static ./static
COPY --from=builder /src/templates ./templates

EXPOSE 8090

# Docker checks this URL every 30s and marks the container unhealthy if it fails
HEALTHCHECK CMD wget -qO- http://localhost:8090/ || exit 1

CMD ["./server"]