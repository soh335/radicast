[![wercker status](https://app.wercker.com/status/6a50ede5efd0abd5b078ab080fe7b61c/s/master "wercker status")](https://app.wercker.com/project/bykey/6a50ede5efd0abd5b078ab080fe7b61c)
# radicast

* record radiko
* serve rss for podcast

## REQUIRE

* rtmpdump
* swftools
* ffmpeg or avconv
* or docker (see docker section)

## INSTALL

```
$ go get github.com/soh335/radicast
```

## USAGE

### SETUP CONFIG.JSON

```
$ radicast --setup > config.json
```

### EDIT CONFIG.JSON

```
$ vim config.json
$ cat config.json

{
  "FMT": [
    "0 0 17 * * *"
  ]
}
```

cron specification is [here](https://godoc.org/github.com/robfig/cron#hdr-CRON_Expression_Format)

### LAUNCH

```
$ radicast
$ curl 127.0.0.1:3355/rss # podcast rss
```

### RELOAD CONFIG.JSON

* reload config when receive HUP signal

## DOCKER

```
$ mkdir workspace
$ cd workspace
$ docker pull soh335/radicast
$ docker run --rm soh335/radicast:latest --setup > config.json
$ docker run --rm -p 3355:3355 -v `pwd`:/workspace soh335/radicast:latest --config /workspace/config.json --output /workspace/output
```

## SEE ALSO

* [ripdiko](https://github.com/miyagawa/ripdiko)
