FROM golang

ARG app_env
ENV APP_ENV $app_env

ARG app_port
ENV APP_PORT $app_port

WORKDIR /go/src/github.com/cyub/ip-scout
COPY . .


RUN go get -d -v ./...
RUN go install -v ./... 

CMD ["ip-scout"]