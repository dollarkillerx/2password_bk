FROM alpine:latest
WORKDIR /app
COPY ./cmd/main/new_project /app
COPY ./configs/ /app/configs
COPY ./static/ /app/static
CMD ["/app/new_project"]
