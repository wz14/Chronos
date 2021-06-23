FROM golang:1.16-buster

ARG Name=AAC

WORKDIR /$Name

ADD go.mod .
ADD go.sum .

RUN go env -w GOPROXY=https://goproxy.cn,direct

RUN go mod download

COPY . /$Name

ENV ID=0

RUN make test

# ENTRYPOINT ["./build/start"]