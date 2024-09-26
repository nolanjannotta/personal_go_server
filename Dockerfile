# syntax=docker/dockerfile:1

FROM golang:latest AS builder

WORKDIR /build


COPY . . 



RUN go mod download
RUN GOOS=linux go build -o ./cmd/personal_server/personal_go_server ./cmd/personal_server


FROM gcr.io/distroless/base-debian12


WORKDIR /app
COPY --from=builder /build/cmd/personal_server/personal_go_server ./personal_go_server
COPY --from=builder /build/.env ./.env  


# COPY --from=builder /build/web/markdown ./web/markdown
# COPY --from=builder /build/web/static ./web/static



EXPOSE 8080
EXPOSE 23234





CMD ["/app/personal_go_server"]


