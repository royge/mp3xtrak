package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

// Scan list the files from directory.
func Scan(dir string, c chan string) error {
	err := filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			c <- path.Join(dir, info.Name())
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// Extract get mp3 audio for video file.
func Extract(c chan string, command, dir string) error {
	for s := range c {
		if s != "" {
			ext := filepath.Ext(s)
			name := filepath.Base(s)

			if ext != "" {
				name = strings.Replace(name, ext, ".mp3", 1)
			}

			out := path.Join(
				dir,
				name,
			)
			cmd := exec.Command(command)
			cmd.Args = []string{"-i", s, out}

			if err := cmd.Run(); err != nil {
				return err
			}

			cmd.Process.Kill()
		}
	}

	return nil
}

func main() {
	s := flag.String("s", "", "Source directory.")
	o := flag.String("o", "", "Output directory.")
	c := make(chan string)

	flag.Parse()

	go func() {
		if err := Scan(*s, c); err != nil {
			log.Fatalf("error scanning directory: %v", err)
		}
		close(c)
	}()

	if err := Extract(c, "ffmpeg", *o); err != nil {
		log.Fatalf("error extracting audio: %v", err)
	}
}
