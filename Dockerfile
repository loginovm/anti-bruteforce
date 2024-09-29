FROM golang:1.22 as build

WORKDIR "/go/src/"

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o bin/anti-bruteforce-app cmd/antibforce/*.go

FROM alpine:3.9

WORKDIR "/opt/anti-bruteforce"

COPY --from=build /go/src/bin/anti-bruteforce-app .
COPY ./migrations ./migrations
COPY ./cmd/antibforce/config.toml .

EXPOSE 8080
ENTRYPOINT [ "./anti-bruteforce-app" ]


