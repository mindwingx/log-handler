FROM golang:1.19

WORKDIR /app

COPY . .

ENV GO111MODULE on
ENV GOPROXY https://goproxy.cn

RUN apt update && apt install nano git

# RUN go test -v .

RUN mv .env.shared .env

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o logger service/main.go
RUN chmod +x logger

#CMD ["/app/main"]