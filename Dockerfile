FROM golang:1.22-alpine as builder
RUN apk --update add build-base

WORKDIR /src/app
ADD go.mod .
RUN go mod download

ADD . .
RUN go build -o bin/db ./cmd/db
RUN go run ./cmd/build

FROM alpine

# Install necessary packages
RUN apk add --no-cache tzdata ca-certificates 
# tor

WORKDIR /bin/

COPY --from=builder /src/app/bin/app .
COPY --from=builder /src/app/bin/db .

# Expose Tor's SOCKS proxy port
# EXPOSE 3000 
# 9050

# Run Tor and then your application
# tor &&
CMD /bin/app
