FROM golang:1.19-alpine AS build

WORKDIR /app

COPY . .
RUN go mod download

RUN go build -v -o /camtarr

FROM alpine
RUN apk add ca-certificates

COPY --from=build /camtarr /camtarr

EXPOSE 9031
CMD ["/camtarr"]
