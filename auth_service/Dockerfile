FROM golang

WORKDIR /app

COPY . .

RUN go mod tidy

EXPOSE 10001

CMD ["go", "run", "cmd/app/main.go"]
