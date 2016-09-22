# mp3xtrak
Extract mp3 audio from mp4 files.

Pre-requisites
--------------
This tool currently uses ffmpeg to do the actual work.
```
sudo apt-get install ffmpeg
```
Installation
------------
```
git clone https://github.com/r0y3/mp3xtrak.git
cd mp3xtrak
```

Configuration
-------------
Modify settings.py and set the following accordingly.
```python
MP4_DIR='<your mp4 directory>'
MP3_DIR='<your target directory for mp3 files>'
```

Run
---
Then execute the main.py
```
./main.py
```
