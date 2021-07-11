FROM golang:1.16-buster as builder
WORKDIR /app/kubenchctl
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o /kubenchctl .
RUN chmod +x /kubenchctl

FROM ubuntu:18.04
COPY --from=builder /kubenchctl /kubenchctl
ENTRYPOINT [ "/kubenchctl" ]
