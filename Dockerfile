FROM golang:1.27-alpine AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o /out/storeapi ./src/cmd/storeapi

FROM alpine:3.22
RUN adduser -D app
WORKDIR /app
COPY --from=build /out/storeapi /app/storeapi
# COPY nginx/static/docs /app/
# COPY migrations /app/migrations
USER app
ENTRYPOINT ["/app/storeapi"]
