FROM golang:alpine AS build
LABEL maintainer="Sam McGeown <smcgeown@vmware.com>"
WORKDIR /go/src/github.com/sammcgeown/cs-cli/
COPY . /go/src/github.com/sammcgeown/cs-cli/
RUN go build


FROM alpine:latest
COPY --from=build /go/src/github.com/sammcgeown/cs-cli/cs-cli .
ENTRYPOINT ["./cs-cli"]