FROM library/golang:1.15.6-alpine3.12 as build
LABEL maintainer="ssheitmann@gmail.com"

ENV GO111MODULE=on

WORKDIR /app
COPY . .
RUN go mod download \
    && CGO_ENABLED=0 \
    GOOS=linux GOARCH=amd64 \
    go build \
    -o ./out/ofever


FROM library/alpine:3.11.6 as runtime
LABEL maintainer="ssheitmann@gmail.com"

RUN addgroup application-group --gid 1001 \
    && adduser application-user --uid 1001 \
    --ingroup application-group \
    --disabled-password

RUN apk add \
    --no-cache \
    ca-certificates \
    apache2-utils

WORKDIR /app
COPY --from=build /app/out .
RUN chown --recursive application-user .
USER application-user
ENTRYPOINT ["/app/ofever"]