# Step 1: Build the Go application
FROM golang:1.24.2 AS builder

WORKDIR /app

# Copy dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the application
RUN go build -o main .

# Step 2: Install GLIBC 2.34 on a custom image
FROM ubuntu:22.04 AS glibc-installer

# Install dependencies for building GLIBC
RUN apt-get update && apt-get install -y build-essential wget

# Download and build GLIBC 2.34
RUN wget https://ftp.gnu.org/gnu/libc/glibc-2.34.tar.gz && \
    tar -xvzf glibc-2.34.tar.gz && \
    cd glibc-2.34 && \
    mkdir build && cd build && \
    ../configure --prefix=/usr && \
    make -j$(nproc) && \
    make install

# Step 3: Prepare a minimal distroless image and copy the application
FROM gcr.io/distroless/base-debian10

# Copy the compiled Go application from the builder stage
COPY --from=builder /app/main /app/main

# Copy the docs and environment file
COPY ./docs /app/docs
COPY .env /app/env

# Copy the GLIBC 2.34 libraries into the distroless image
COPY --from=glibc-installer /usr /usr

# Set the working directory
WORKDIR /app

# Expose the port the app will run on
EXPOSE 8080

# Command to run the application
CMD ["/app/main"]
