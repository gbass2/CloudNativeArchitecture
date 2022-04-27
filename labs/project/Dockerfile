FROM golang:1.17-alpine AS build

WORKDIR /go/src/gamepub/

ADD . /go/src/gamepub/

RUN go mod download -x

RUN CGO_ENABLED=0 go build -o /bin/gamepub

EXPOSE 50051

FROM scratch
COPY --from=build /bin/gamepub /bin/gamepub
ENTRYPOINT ["/bin/gamepub"]
