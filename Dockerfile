# use goland:lastest instead of golang:alpine because go git is not available in alpine version
FROM golang:1.20 as builder
WORKDIR /go/src/ngrok-file-server
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 go build -o bin/api -v api/main.go

FROM scratch
COPY --from=builder /go/src/ngrok-file-server/bin/api /app/bin/
ENTRYPOINT ["/app/bin/api"]