FROM golang:1.17-alpine AS build

WORKDIR /go/src/movieserver/

ADD . /go/src/movieserver/

RUN go mod download -x

RUN CGO_ENABLED=0 go build -o /bin/movieserver

EXPOSE 50051

FROM scratch
COPY --from=build /bin/movieserver /bin/movieserver
ENTRYPOINT ["/bin/movieserver"]
