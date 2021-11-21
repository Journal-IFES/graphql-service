FROM alpine:edge

ARG COMMIT

COPY "./graphql-service-${COMMIT}" /app
ENTRYPOINT [ "/app/graphql-service-${COMMIT}" ]