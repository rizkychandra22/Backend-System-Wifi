FROM golang:1.25-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /wifi-api .

FROM alpine:3.21
RUN apk add --no-cache tzdata ca-certificates
COPY --from=build /wifi-api /wifi-api
EXPOSE 8080
ENTRYPOINT ["/wifi-api"]
