FROM golang

WORKDIR /app

COPY . .

RUN go mod tidy

#EXPOSE 10003

CMD ["go", "run", "cmd/app/main.go"]
