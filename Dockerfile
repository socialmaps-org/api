FROM golang:1.26.0-trixie

WORKDIR /usr/src/app

# Pre-copy/cache go.mod for pre-downloading dependencies and re-downloading them
# in subsequent builds only if they change.
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN make

FROM debian:trixie

WORKDIR /root

COPY --from=0 /usr/src/app/bin/socialmaps-api /bin/socialmaps-api
ENTRYPOINT ["/bin/socialmaps-api"]
