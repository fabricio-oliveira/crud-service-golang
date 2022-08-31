FROM golang:1.18-alpine3.16 AS builder

ENV GO111MODULE=on 

WORKDIR /app
 
COPY .. .

ENV TZ America/Chicago

#dependencies
RUN apk add --no-cache make  build-base
RUN go mod tidy 

WORKDIR /app/cmd/api
 
RUN GOOS=linux go build -a -installsuffix cgo -o app .


#final image
FROM alpine:3.16
WORKDIR /root/
COPY --from=builder /app/cmd/api/app .
CMD ["./app"]