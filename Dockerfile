FROM golang:1.18 as build

WORKDIR /build
COPY . ./
RUN make

FROM ubuntu

WORKDIR /opt
COPY --from=build /build/minic ./

USER nobody
ENTRYPOINT ["/opt/minic"]
