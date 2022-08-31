FROM golang:1.18-alpine3.16 AS builder

ENV GO111MODULE=on 

WORKDIR /app
 
COPY .. .

ENV TZ America/Chicago

#dependencies
RUN go mod tidy 

WORKDIR /app/cmd/crudservice
 
RUN GOOS=linux go build -a -installsuffix cgo -o app .


#Second stage for a smaller image
FROM alpine:3.16
WORKDIR /root/
COPY --from=builder /app/cmd/tutorial/app .
CMD ["./app"]