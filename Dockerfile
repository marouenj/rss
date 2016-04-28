FROM alpine:3.3

MAINTAINER marouenj

RUN mkdir /app
RUN mkdir /files

COPY ./rss /app/rss
ENV PATH $PATH:/app

WORKDIR /files

ENTRYPOINT ["rss"]
CMD ["--base_dir=./"]
