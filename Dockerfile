FROM alpine:edge

ARG COMMIT

RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
COPY "./graphql-service-${COMMIT}" "/app/graphql-service-${COMMIT}"
ENTRYPOINT [ "/app/graphql-service-${COMMIT}" ]