FROM golang:latest
 
RUN mkdir -p /app
 
WORKDIR /app
ADD ./src /app
 
RUN rm -rf go.mod go.sum
RUN go mod init gifcreator
RUN go mod tidy
RUN go build -o app .
 
CMD ["./app"]