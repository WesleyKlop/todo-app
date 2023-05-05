FROM golang:1.20
WORKDIR /app
ADD bin bin
ENTRYPOINT ['/app/bin/todo-api']