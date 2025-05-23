FROM golang:1.24-alpine AS build

WORKDIR /work

COPY cmd cmd
COPY internal internal
COPY go.mod go.sum ./
RUN --mount=type=cache,id=gomod,target=/go/pkg/mod \
    --mount=type=cache,id=gobuild,target=/root/.cache/go-build \
    go mod download

ENV CGO_ENABLED=0
#ENV GOOS=darwin
#ENV GOARCH=arm64

RUN go build -o /tmp ./cmd/server

FROM scratch AS  deploy
COPY --from=build /tmp/server /server
EXPOSE 8080
CMD ["/server"]

