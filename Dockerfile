FROM golang:1.17-alpine

ADD ./graphql-service /app
ENTRYPOINT [ "/app/graphql-service" ]