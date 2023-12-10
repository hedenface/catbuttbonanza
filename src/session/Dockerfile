FROM debian:bookworm-slim

EXPOSE 8081/tcp

COPY resources ./resources
COPY session .

ENTRYPOINT ["./session"]
