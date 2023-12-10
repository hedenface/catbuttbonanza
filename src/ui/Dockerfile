FROM debian:bookworm-slim

EXPOSE 8080/tcp

COPY resources ./resources
COPY ui .

ENTRYPOINT ["./ui"]
