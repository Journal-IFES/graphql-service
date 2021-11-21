FROM golang:1.17-alpine

COPY ./graphql-service /app
ENTRYPOINT [ "/app/graphql-service" ]