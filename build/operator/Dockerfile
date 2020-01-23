FROM alpine:3.9

RUN mkdir /app
WORKDIR /app

COPY _output/operator operator
RUN chmod +x operator

CMD ["/app/operator"]

