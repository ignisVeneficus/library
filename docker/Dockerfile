# Build stage
FROM golang:1.23-alpine AS builder

RUN apk --no-cache add git

WORKDIR /usr/local/bin/library

RUN git clone --depth 1 -b main https://github.com/ignisVeneficus/library .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /library/library .

FROM alpine:latest

RUN apk --no-cache add ca-certificates shadow tzdata su-exec

ENV IMAGE_DIR=/covers
ENV BOOKS_DIR=/books
ENV PUID=1000
ENV PGID=1000


RUN mkdir -p ${IMAGE_DIR} ${BOOKS_DIR}

COPY --from=builder /library/library /usr/local/bin/library
COPY --from=builder /usr/local/bin/library/web ./web


ENTRYPOINT ["/bin/sh", "-c", "\
  if ! getent group ${PGID} >/dev/null 2>&1; then \
    addgroup -g ${PGID} library; \
  fi && \
  # User létrehozása, ha nem létezik \
  if ! id -u ${PUID} >/dev/null 2>&1; then \
    adduser -D -H -u ${PUID} -G library library; \
  fi && \
  if [ -n \"${TZ}\" ]; then \
    ln -snf /usr/share/zoneinfo/${TZ} /etc/localtime && \
    echo ${TZ} > /etc/timezone ; \
    echo \"Timezone set to ${TZ}\" ; \
  else \
    echo \"No TZ variable provided, using default (UTC)\" ; \
  fi && \
  chown -R library:library ${IMAGE_DIR} ${BOOKS_DIR} && \
  export LIBRARY_BOOKS=${BOOKS_DIR} && \
  export LIBRARY_COVERS=${IMAGE_DIR} && \
  echo 'User setup done. Starting application as library...' && \
  exec su-exec library /usr/local/bin/library \
"]

EXPOSE 8888