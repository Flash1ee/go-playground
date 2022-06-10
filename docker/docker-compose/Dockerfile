FROM alpine:latest

RUN apk update && apk add curl

RUN apk add bash

COPY curl.sh .

CMD ["bash", "./curl.sh"]