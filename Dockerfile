FROM golang:1.24.4-bullseye

RUN apt-get update -y && apt install git inotify-tools -y

ENV CGO_ENABLED=0
ENV GO111MODULE=on

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
WORKDIR $GOPATH/src/confluence_cli

COPY . .

RUN go mod download
RUN GOOS=linux go build -o confluence_cli
RUN chmod +x confluence_cli
RUN mv confluence_cli /usr/local/bin/confluence_cli
ENTRYPOINT ["confluence_cli"]