FROM golang:1.21 AS build-stage 
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/app main.go

FROM docker
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
COPY --from=build-stage /go/bin/app /go/bin/app
ENTRYPOINT [ "/go/bin/app" ]
