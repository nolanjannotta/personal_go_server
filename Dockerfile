# syntax=docker/dockerfile:1

FROM golang:latest AS builder

WORKDIR /build


COPY . . 


RUN go mod download
RUN GOOS=linux go build -o ./personal_go_server


FROM gcr.io/distroless/base-debian12


WORKDIR /app
COPY --from=builder /build/personal_go_server ./personal_go_server
COPY --from=builder /build/markdown ./markdown
COPY --from=builder /build/static ./static


EXPOSE 8080
EXPOSE 23234





CMD ["/app/personal_go_server"]


