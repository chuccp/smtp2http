# syntax=docker/dockerfile:1

FROM golang:1.23

# Set destination for COPY
WORKDIR /app
# Download Go modules
COPY go.mod go.sum ./
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY="https://goproxy.cn"
RUN go mod download


# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY . ./


##go build -o C:/Users/cooge/software/httpPush/httpPush.exe github.com/chuccp/httpPush
# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /smtp2http .


ADD https://github.com/chuccp/d-mail-view/releases/latest/download/d-mail-view.zip /

CMD [ "/smtp2http -unzip d-mail-view.zip web" ]

# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can (optionally) document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/engine/reference/builder/#expose
EXPOSE 12566 12567

# Run
CMD [ "/smtp2http -web_port 12566 -api_port 12577" ]
