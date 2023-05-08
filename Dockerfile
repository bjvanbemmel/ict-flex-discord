FROM golang:1.20.4-alpine3.17

WORKDIR /ict-flex-discord

COPY ./ ./

RUN go mod download

RUN go build -o /bot .

CMD [ "/bot" ]
