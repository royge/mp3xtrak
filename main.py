#!/usr/bin/env python3

import os
import multiprocessing
import subprocess
import glob
import re
import argparse

from threading import Thread
from queue import Queue

from settings import MP4_DIR, MP3_DIR

class Extractor():
    def __init__(self, queue, dest_dir):
        self.queue = queue
        self.dest_dir = dest_dir

    def extract(self):
        file_name = re.escape(self.queue.get())
        mp3_file_name = os.path.basename(file_name.replace('.mp4', '.mp3'))
        print("Extracting from", file_name.replace("\\", ""), "...")
        command = "ffmpeg -i {} {}/{}".format(
            file_name,
            self.dest_dir,
            mp3_file_name
        )
        res = subprocess.Popen(
            command,
            shell=True,
            stdin=subprocess.PIPE,
            stdout=subprocess.PIPE,
            stderr=subprocess.STDOUT
        ).communicate()

        print("Done extracting {}.".format(mp3_file_name.replace("\\", "")))

class Worker(Thread):
    def __init__(self, extractor):
        Thread.__init__(self)
        self.extractor = extractor

    def run(self):
        while True:
            self.extractor.extract()
            self.extractor.queue.task_done()


class Scanner():
    def __init__(self, queue):
        self.queue = queue

    def scan(self, base_dir):
        pattern = ["*.mp4"]
        for root, dirs, files in os.walk(base_dir):
            for trace in pattern:
                pattern_path = os.path.join(root, trace)
                for filename in glob.glob(pattern_path):
                    if os.path.exists(filename):
                        self.queue.put(filename)


def main():
    parser = argparse.ArgumentParser(description='Extract mp3 from mp4 videos.')
    parser.add_argument('--mp4dir', default=MP4_DIR, help='MP4 directory')
    parser.add_argument('--mp3dir', default=MP3_DIR, help='MP3 directory')

    args = parser.parse_args()

    queue = Queue()

    scanner = Scanner(queue)
    scanner.scan(args.mp4dir)

    for i in range(multiprocessing.cpu_count()):
        worker = Worker(Extractor(queue, args.mp3dir))
        worker.daemon = True
        worker.start()

    queue.join()

if __name__ == '__main__':
    main()
