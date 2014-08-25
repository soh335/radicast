# radicast

* record radiko
* serve rss for podcast

## REQUIRE

* rtmpdump
* swftools
* ffmpeg

## INSTALL

```
$ go get github.com/soh335/radicast
```

## USAGE

### SETUP CONFIG.JSON

```
$ radicast --setup # create config.json
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

### LAUNCH

```
$ radicast
$ curl 127.0.0.1:3355/rss # podcast rss
```

## SEE ALSO

* [ripdiko](https://github.com/miyagawa/ripdiko)
