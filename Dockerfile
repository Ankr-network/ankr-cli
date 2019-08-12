FROM golang:1.12-alpine as builder
RUN apk update && \
    apk add git && \
    apk add --update bash && \
    apk add openssh && \
    apk add make

COPY id_rsa /root/.ssh/
RUN ssh-keyscan github.com >> ~/.ssh/known_hosts
RUN chmod go-w /root
RUN chmod 700 /root/.ssh
RUN chmod 600 /root/.ssh/id_rsa
RUN export GO111MODULE=on

WORKDIR $GOPATH/src/github.com/Ankr-network/ankr-cli/
COPY . $GOPATH/src/github.com/Ankr-network/ankr-cli/

RUN go mod download
RUN CGO_ENABLED=0 \
        GOOS=linux \
        GOARCH=amd64 \
        go build -a \
        -installsuffix cgo \
        -o /go/bin/ankr-cli \
        main.go

FROM alpine:3.7
RUN  apk update && \
     apk add libc6-compat && \
     apk add ca-certificates
COPY --from=builder /go/bin/ankr-cli /bin/ankr-cli
