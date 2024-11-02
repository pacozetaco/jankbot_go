
FROM golang:1.23.2-alpine3.20

# Set the working directory in the container to /jankbot2
WORKDIR /jankbot_go

# Copy the current directory contents into the container at /jankbot
COPY . .

# Install system dependencies and Python packages
RUN apk update && \
    apk add --no-cache ffmpeg && \
    rm -rf /var/cache/apk/*

# Run main.py when the container launches
CMD ["go", "run main.go"]
