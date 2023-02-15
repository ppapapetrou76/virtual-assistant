FROM golang:1.20.1-alpine

LABEL "com.github.actions.name"="Issue and Pull Request virtual assistant"
LABEL "com.github.actions.description"="Automates tasks on pull requests and issues actions"
LABEL "com.github.actions.icon"="award"
LABEL "com.github.actions.color"="blue"
LABEL "maintainer"="Patroklos Papapetrou <ppapapetrou76@gmail.com>"
LABEL "repository"="https://github.com/ppapapetrou76/virtual-assistant"

RUN apk add --no-cache git

WORKDIR /go/src/app
COPY . .
ENV GO111MODULE=on
RUN go build -o action ./cmd
ENTRYPOINT ["/go/src/app/action"]
