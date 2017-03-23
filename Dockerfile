FROM debian:jessie

COPY cg-fake-uaa /bin

EXPOSE 8080

CMD /bin/cg-fake-uaa
