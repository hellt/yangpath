FROM scratch

LABEL maintainer="Roman Dodin <dodin.roman@gmail.com>"
LABEL documentation="https://yangpath.netdevops.me"
LABEL repo="https://github.com/hellt/yangpath"

COPY yangpath /app/yangpath
ENTRYPOINT [ "/app/yangpath" ]
CMD [ "version" ]
