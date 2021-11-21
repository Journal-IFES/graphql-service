FROM alpine:edge

ARG COMMIT

RUN apk update
RUN apk add --no-cache libc6-compat

COPY "./graphql-service-$COMMIT" "/app/graphql-service-$COMMIT"
ENTRYPOINT [ "/app/graphql-service-$COMMIT" ]