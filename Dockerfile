# Stage 1: Build Go application
FROM golang:1.17 AS build



WORKDIR /app

# Copy Go source files
COPY . .

# Grant execution permissions
RUN chmod +x /app/start.sh

# Build the Go application
RUN go build -o main .

# Stage 2: Create final Docker image
FROM mysql:latest

# Set MySQL environment variables
ENV MYSQL_ROOT_PASSWORD=admin
ENV MYSQL_DATABASE=synapso

# Copy the built Go binary from the first stage
COPY --from=build /app/main /app/main
COPY start.sh /app/


# Copy SQL initialization script
COPY init.sql /docker-entrypoint-initdb.d/

# Expose the MySQL port
EXPOSE 3306
EXPOSE 8080

# Start MySQL server
CMD ["/app/start.sh"]

