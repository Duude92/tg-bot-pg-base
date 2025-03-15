FROM golang:1.24
LABEL authors="Duude92"
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -v -o /app

CMD ["/app"]