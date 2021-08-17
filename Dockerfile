FROM golang:1.15

WORKDIR /usr/src/app

COPY . .

RUN go get .
RUN go build -o app .

EXPOSE 3000

CMD [ "./app" ]