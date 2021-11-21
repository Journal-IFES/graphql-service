FROM golang:1.17-alpine

ARG COMMIT

COPY "./graphql-service-${COMMIT}" /app
ENTRYPOINT [ "/app/graphql-service-${COMMIT}" ]