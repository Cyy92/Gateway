# Build fxwatcher
FROM golang:1.19

ENV GO111MODULE=on

RUN mkdir -p ${GOPATH}/src/github.com/keti-openfx/openfx/executor/go
WORKDIR ${GOPATH}/src/github.com/keti-openfx/openfx/executor/go

COPY . .

RUN go build -o fxwatcher .

