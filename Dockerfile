FROM alpine:3.12
RUN apk --no-cache add ca-certificates
WORKDIR /zt100
COPY ./local/build/zt100 /usr/local/bin/zt100
COPY ./tmpl /zt100/tmpl
EXPOSE 8080
CMD ["/usr/local/bin/zt100"]