# Use a lightweight Go image with Alpine Linux
FROM golang:1.23.2-alpine3.20

# Set the working directory in the container to /jankbot_go
WORKDIR /jankbot_go

# Copy the current directory contents into the container at /jankbot_go
COPY . .

# Install system dependencies (e.g., ffmpeg)
RUN apk update && \
    apk add --no-cache ffmpeg && \
    rm -rf /var/cache/apk/*

# Run main.go when the container launches
CMD ["go", "run", "/jankbot_go/main.go"]
