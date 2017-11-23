FROM golang
RUN go get -u -v github.com/myml/ytbdl
WORKDIR /go/src/github.com/myml/ytbdl/
CMD go run /go/src/github.com/myml/ytbdl/main.go
