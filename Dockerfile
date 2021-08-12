FROM golang:latest

LABEL maintainer="Joking Devon"

WORKDIR /app

COPY go.mod .

COPY go.sum .

RUN go mod download

COPY . .

ENV MONGO_USER="<you mongo user>"

ENV MONGO_PASS="<your mongo user's password>"

ENV MONGO_DB="<your mongo database name>"

RUN go build

EXPOSE 12345

CMD ./task-planner