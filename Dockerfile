FROM golang:1.9.2 as builder
WORKDIR /go/src/github.com/samkreter/acictl/ 
COPY . /go/src/github.com/samkreter/acictl/
# RUN go test ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o run .

FROM alpine:latest  
RUN apk --update add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/samkreter/acictl/run .
CMD ["./run"]