# syntax=docker/dockerfile:1

FROM scratch

WORKDIR /

COPY bin/workspaceone-exporter /

EXPOSE 9740

ENTRYPOINT ["/workspaceone-exporter"]
