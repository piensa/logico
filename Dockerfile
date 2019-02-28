FROM golang:latest
RUN mkdir /app 
ADD . /app/
RUN go build -o main .
WORKDIR /app 
CMD ["/app/main"]