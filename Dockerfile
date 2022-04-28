FROM golang:1.16-alpine as builder
COPY . $GOPATH/work/
WORKDIR $GOPATH/work/
RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/go-landing

FROM scratch
COPY --from=builder /go/bin/go-landing /go/bin/go-landing
EXPOSE 9000
ENTRYPOINT ["/go/bin/go-landing"]
