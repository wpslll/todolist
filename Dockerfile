FROM golang:1.25.8-bookworm

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o /app/exe main.go

CMD [ "/app/exe" ]