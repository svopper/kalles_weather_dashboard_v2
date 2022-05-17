FROM golang:1.18-alpine

ENV DMI_API_KEY=xxx
ENV GIN_MODE=release

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /app

EXPOSE 8080

CMD [ "/app" ]