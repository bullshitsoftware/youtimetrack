FROM golang:1.18.3-alpine

RUN apk add gcc musl-dev curl git openssh fish starship

RUN go install github.com/cweill/gotests/gotests@latest
RUN go install github.com/fatih/gomodifytags@latest
RUN go install github.com/josharian/impl@latest
RUN go install github.com/haya14busa/goplay/cmd/goplay@latest
RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN go install honnef.co/go/tools/cmd/staticcheck@latest
RUN go install golang.org/x/tools/gopls@latest

RUN adduser -u 1000 -D -s /usr/bin/fish user
USER user

RUN mkdir -p ~/.config/fish && echo "starship init fish | source" > ~/.config/fish/config.fish
