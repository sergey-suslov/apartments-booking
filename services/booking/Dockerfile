FROM golang:1.14

WORKDIR /go/src/app
COPY . .


RUN go get -d -v ./...
RUN go install -v ./...
RUN go build cmd/main.go


FROM golang:1.14

WORKDIR /go/src/app
COPY --from=0 /go/src/app/main .

CMD ["./main"]