FROM golang:1.17

ADD image-updater /image-updater

ENTRYPOINT ["/image-updater"]
