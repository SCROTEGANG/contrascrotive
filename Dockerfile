FROM golang:1.17 AS builder

WORKDIR /contrascrotive

COPY ./go.mod ./go.sum

RUN go mod download

COPY ./ ./

RUN go build . -o contrascrotive

FROM gcr.io/distroless/base:nonroot
COPY --from=builder /contrascrotive/contrascrotive /contrascrotive
ENTRYPOINT [ "/contrascrotive" ]