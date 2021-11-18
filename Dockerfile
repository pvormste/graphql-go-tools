# syntax=docker/dockerfile:1

FROM golang:1.15-alpine
WORKDIR /graphql-go-tools

#COPY go.mod ./
#COPY go.sum ./
COPY . .

RUN go mod download

CMD ["sh", "-c", "cd examples/federation ; sh start.sh"]
