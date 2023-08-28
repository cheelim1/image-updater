FROM golang:1.23

ADD image-updater /image-updater

ENTRYPOINT ["/image-updater"]
