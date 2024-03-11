# syntax=docker/dockerfile:1
#Build Stage
FROM golang:1.19-alpine AS builder

WORKDIR /src

COPY . ./
RUN apk add build-base && go build -o forum ./cmd/main/main.go
# Deploy Stage
FROM alpine
WORKDIR /src
LABEL name="forum" \
    maintainer="Semyon Serbulov <sonnenblumenglas@gmail.com>" \
    org.label-schema.schema-version="1.0" \
    org.label-schema.name="forum" \
    org.label-schema.vcs-url="https://01.alem.school/git/sserbulo/forum" \
    org.label-schema.docker.cmd="docker run -p 8000:8000 -d --name forum-container forum"
COPY --from=builder /src .
EXPOSE 3000
ENTRYPOINT ["./forum"]