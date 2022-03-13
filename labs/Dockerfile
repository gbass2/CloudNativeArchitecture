FROM golang:1.17-alpine AS build

#ENV GO111MODULE=on

COPY ./lab7/movieapi /go/src/

WORKDIR /go/src/labs/
COPY ./lab7 .
#COPY ./lab7/movieapi ./labs
COPY go.mod .
COPY go.sum .

RUN go mod download -x
#RUN go get ./lab7/movieapi
#RUN go install ./lab7/movieapi

RUN CGO_ENABLED=0 go build -o /bin/labs

FROM scratch
COPY --from=build /bin/labs /bin/labs
ENTRYPOINT ["/bin/labs"]
