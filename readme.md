# GoAudiobook

GoAudiobook is a lightweight app

## Install

Local Install (using go tool)

```
go install -g github.com/zedawg/goaudiobook
goaudiobook -port 8080 -dir ~/downloads/audiobooks
```

Using Docker (using docker)

```
docker pull dockerhub.com/zedawg/goaudiobook
docker run -v data:/data -p 80:8000 goaudiobook:latest
```
