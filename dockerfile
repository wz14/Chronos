FROM golang:1.16-buster
WORKDIR /AAC
ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . /AAC