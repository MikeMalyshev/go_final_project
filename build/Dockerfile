# Explain of this https://github.com/confluentinc/confluent-kafka-go/issues/898
FROM --platform=linux/amd64 golang:1.22.2

WORKDIR /app
COPY go.sum go.mod ./
RUN go mod download

COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY web/ ./web/
COPY bin ./bin/

ENV CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

RUN go build -C cmd/ -o /bin/todo_app

ENV TODO_PORT="7540" \
    TODO_WEBDIR="/app/web" \
    TODO_DBPATH="/app/scheduler.db"

EXPOSE $TODO_PORT

CMD [ "/bin/todo_app" ]
