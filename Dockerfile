FROM scratch

LABEL maintainer = "wayne900619@gmail.com"

EXPOSE 1225

COPY ./Mr.Coding-LineBot /

WORKDIR /

ENTRYPOINT ["./Mr.Coding-LineBot"]