FROM golang:1.22.6-bookworm AS build

COPY . /workspace

WORKDIR /workspace

RUN CGO_ENABLED=0 go build -o serve main.go

FROM scratch

COPY --from=build /workspace/serve /usr/local/bin/

ENTRYPOINT [ "/usr/local/bin/serve" ]
