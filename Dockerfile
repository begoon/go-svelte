FROM scratch

COPY tmp/main /exe

ENTRYPOINT ["/exe"]
