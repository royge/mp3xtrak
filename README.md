# mp3xtrak

Extract mp3 audio from mp4 files.

## Pre-requisites

This tool currently uses [ffmpeg](https://www.ffmpeg.org/) to do the actual work.

**NOTE:** Only tested on *Ubuntu 16.04 LTS* and *MacOS Sierra*

A. Ubuntu

	```
	$ sudo apt-get install ffmpeg
	```

B. Mac

	```
	$ brew install ffmpeg
	```

## Installation

1. Install [go](https://golang.org)

1. Install `mp3xtrak`

	**A. Using `go get`**

	```
	$ go get github.com/royge/mp3xtrak
	```

	**B. From source**

	```
	$ git clone https://github.com/royge/mp3xtrak.git
	$ cd mp3xtrak
	$ go build -o mp3xtrak
	$ go install
	```

## How to Use

```
$ mp3xtrak -s=<your-mp4-directory> -o=<your-music-directory>
```
