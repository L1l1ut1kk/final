FROM golang:alpine as builder


WORKDIR /final
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /final_web

CMD ["/final_web"]
