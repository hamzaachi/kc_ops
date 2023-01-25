#docker run -v /home/hachi/config:/config --rm box.eadn.dz:8888/kc_ops:01 --path=/config/config.yml --roles
FROM golang:alpine

RUN mkdir /src
WORKDIR /src

COPY ./go.mod /src/go.mod
COPY ./go.sum /src/go.sum
RUN go mod download

COPY ./ /src/
RUN go build -o kc_ops .

FROM alpine:latest
RUN apk --no-cache add \
    ca-certificates \
    git
    
RUN   addgroup -g 70 -S kc_ops; \
	adduser -u 70 -S -D -G kc_ops -H -h /app -s /bin/sh kc_ops;

RUN git config --global --add safe.directory '*'

WORKDIR /app/

COPY --chown=kc_ops --from=0 /src/kc_ops ./

USER kc_ops

ENTRYPOINT ["/app/kc_ops"]
