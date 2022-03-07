FROM golang:1.17 AS build
WORKDIR /go/src/github.com/alias-rahil/cluster-configs
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o entrypoint

FROM scratch
COPY --from=build /go/src/github.com/alias-rahil/cluster-configs/entrypoint entrypoint

ENTRYPOINT [ "./entrypoint" ]
