FROM golang as builder
WORKDIR /src/
COPY main.go go.mod go.sum /src/
RUN CGO_ENABLED=0 go build

FROM alpine
COPY --from=builder /src/vw-oauth /
ENTRYPOINT ["/vw-oauth"]
