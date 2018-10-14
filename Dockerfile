FROM golang:1.11.1-alpine3.8 As builder

WORKDIR /go/src/github.com/kiririmode/grpc-sandbox
COPY . .

RUN apk --no-cache add protobuf make git
RUN go get -u github.com/golang/protobuf/protoc-gen-go
RUN make pb deps
RUN GOOS=linux go build -o grpc-server .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /go/src/github.com/kiririmode/grpc-sandbox/grpc-server ./
EXPOSE 8000
CMD ["/app/grpc-server"]
