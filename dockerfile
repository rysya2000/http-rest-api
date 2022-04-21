FROM golang

WORKDIR /web

COPY . .

RUN go build -v ./cmd/apiserver

CMD ["./apiserver"]