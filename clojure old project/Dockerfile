FROM clojure:lein-2.8.1
ADD . /app
WORKDIR /app
EXPOSE 18083
CMD ["lein", "ring", "server"]
