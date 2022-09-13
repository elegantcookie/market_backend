FROM golang

WORKDIR /app

COPY . .

RUN go mod tidy

EXPOSE 10002

CMD ["go", "run", "cmd/app/main.go"]
