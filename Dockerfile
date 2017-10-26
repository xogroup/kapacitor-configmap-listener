FROM golang:latest

MAINTAINER Lam Chan (lchan@xogrp.com)

ENTRYPOINT [ "kapacitor-configmap-listener" ]

CMD [ "-h" ]

COPY kapacitor-configmap-listener /usr/local/bin/


