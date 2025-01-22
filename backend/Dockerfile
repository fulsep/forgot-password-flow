FROM golang AS build

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o backend main.go

# RUN go build -mod vendor -o backend main.go

FROM debian:bookworm-slim

COPY --from=build /app/backend /usr/local/bin/backend

RUN chmod +x /usr/local/bin/backend

WORKDIR /app

CMD [ "backend" ]