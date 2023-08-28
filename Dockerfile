FROM golang:1.17

ADD update-image-tag /update-image-tag

ENTRYPOINT ["/update-image-tag"]
