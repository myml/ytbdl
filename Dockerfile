FROM golang
RUN go get -u -v github.com/myml/ytbdl
EXPOSE 4000
CMD go run /go/src/github.com/myml/ytbdl/main.go
