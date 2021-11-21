FROM alpine:edge

ARG COMMIT

COPY "./graphql-service-${COMMIT}" "/app/graphql-service-${COMMIT}"
ENTRYPOINT [ "/app/graphql-service-${COMMIT}" ]