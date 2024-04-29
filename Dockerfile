# Stage 1: Build the Go application
FROM golang:1.22-alpine as builder

# Install build dependencies
RUN apk --update add build-base

# Set the working directory
WORKDIR /src/app

# Copy the Go module files and download dependencies
COPY go.mod .
RUN go mod download

# Copy the application source code
COPY . .

# Build the application
RUN go build -o /bin/db ./cmd/db
RUN go run ./cmd/build

# Stage 2: Run Tor and the Go application

# FROM builder as inspector

# Set the working directory
# WORKDIR /goapp

# List the contents of the directory
# RUN ls
# RUN ls /bin

FROM dperson/torproxy

# # Set the working directory
WORKDIR /goapp

# # Install necessary packages
RUN apk add --no-cache tzdata ca-certificates tor
ENV PATH=$PATH:/usr/bin/tor

# # # Copy the built binaries from the builder stage
COPY --from=builder /bin/db /goapp/db
COPY --from=builder /bin/app /goapp/app
COPY torrc/torrc "/etc/tor/torrc"

# # # Expose the port used by your Go application
# # EXPOSE 8080

# # # Command to start Tor and then your application
# CMD ["/usr/bin/tor", "-f", "/etc/tor/torrc"] && /goapp/app
CMD /goapp/app
# RUN cat /etc/tor/torrc
