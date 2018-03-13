# mp3xtrak

Extract mp3 audio from mp4 files.

## Pre-requisites

This tool currently uses [ffmpeg](https://www.ffmpeg.org/) to do the actual work.

**NOTE:** Only tested on *Ubuntu 16.04 LTS*
```
sudo apt-get install ffmpeg
```

## Installation

1. Install [go](https://golang.org)

2. Install `mp3xtrak`

```
$ go get github.com/royge/mp3xtrak
```

3. Install from source


```
$ git clone https://github.com/royge/mp3xtrak.git
$ cd mp3xtrak
$ go build -o mp3xtrak
```

## How to Use

```
mp3xtrak -s=<your-mp4-directory> -o=<your-music-directory>
```
