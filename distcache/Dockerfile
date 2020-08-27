 

FROM golang:latest as build
WORKDIR /go/src/distcache
ADD . /go/src/distcache
RUN go get -d -v ./...
RUN go build -o /go/bin/distcache main.go

FROM gcr.io/distroless/base
COPY --from=build /go/bin/distcache /
ENTRYPOINT [ "/distcache" ]