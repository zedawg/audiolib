### Features

- Auto-index audiobooks
- Auto-convert mp3/wav/mpeg to m4b
- Web browser, iOS player

### Install

#### Option #1: go cmd

```
go install -g github.com/zedawg/librarian
librarian [sources...]
```

#### Option #2: Docker

```
docker pull dockerhub.com/zedawg/librarian
docker run -v data:/data -v source:/source -p 8000:8000 librarian:latest
```

