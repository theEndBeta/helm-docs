FROM alpine:3.11.5

COPY yaml-docs /usr/bin/

WORKDIR /yaml-docs

ENTRYPOINT ["yaml-docs"]
