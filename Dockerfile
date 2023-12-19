# syntax=docker/dockerfile:1
# Choose a base image that contains the Go runtime
FROM golang:1.20
# Setting working directory inside the container where app will be placed
WORKDIR /app
# Copy app code into the container
# This copies all the files from current directory into /app directory in the container
COPY . .
# Building an app inside the container
RUN go build -o forumdocker
# Because app listens on a port 8080, here that port is exposed inside the Docker container
EXPOSE 8080
# This runs compiled Go app when the container starts
CMD ["./forumdocker"]
