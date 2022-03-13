FROM golang:1.17-alpine AS build

#ENV GO111MODULE=on

COPY ./movieapi /go/src/

WORKDIR /go/src/movieserver/

ADD . /go/src/movieserver/
#COPY ./lab7/server.go ./

RUN go mod download -x

RUN CGO_ENABLED=0 go build -o /bin/movieserver

EXPOSE 50051

FROM scratch
COPY --from=build /bin/movieserver /bin/movieserver
ENTRYPOINT ["/bin/movieserver"]
