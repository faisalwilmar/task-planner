FROM golang:latest
# use image golang:latest for base image

LABEL maintainer="Joking Devon"

WORKDIR /app
# change working directory

COPY go.mod .

COPY go.sum .

RUN go mod download
# download all dependencies for faster image creation

COPY . .
# cache all the dependencies downloaded

ENV MONGO_USER="<you mongo user>"

ENV MONGO_PASS="<your mongo user's password>"

ENV MONGO_DB="<your mongo database name>"
# set all the environment variables

RUN go build
# build the app

EXPOSE 12345
# expose port for binding

CMD ./task-planner
# run golang application

# Dockerfile reference https://docs.docker.com/get-started/02_our_app/