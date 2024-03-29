FROM golang:1.21.6
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
COPY ./views ./views
COPY .env.docker .
RUN CGO_ENABLED=0 GOOS=linux go build -o /gravityapi
EXPOSE 3000
ENTRYPOINT ["/gravityapi"]