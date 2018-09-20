FROM hub.ibitcome.com/library/golang:1.10-alpine
WORKDIR /app
ADD . /app
RUN mkdir /app/logs && chown -R admin:admin /app/logs
USER admin
#ENV GOPATH="/app"
#RUN go build -o login_server ./src/server/main.go
CMD ["./dgatewayWebBrowser"]
