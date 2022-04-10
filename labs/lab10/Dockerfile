FROM golang:1.17-alpine AS build

WORKDIR /src/
COPY webserver.go /src/
COPY go.mod /src/
COPY go.sum /src/

RUN go mod download -x

RUN CGO_ENABLED=0 go build -o /bin/webserver
FROM scratch
COPY --from=build /bin/webserver /bin/webserver
ENTRYPOINT ["/bin/webserver"]
