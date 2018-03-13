package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

// Scan list the files from directory.
func Scan(dir string, c chan string, exts string) error {
	err := filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			if strings.ContainsAny(info.Name(), exts) {
				c <- path.Join(dir, info.Name())
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// Extract get mp3 audio for video file.
func Extract(s string, command, dir string) error {
	if s != "" {
		s = strings.Replace(s, " ", "\\ ", -1)
		s = strings.Replace(s, "(", "\\(", -1)
		s = strings.Replace(s, ")", "\\)", -1)
		ext := filepath.Ext(s)
		name := filepath.Base(s)

		if ext != "" {
			name = strings.Replace(name, ext, ".mp3", 1)
		}

		p := path.Join(
			dir,
			name,
		)
		cmd := exec.Command(command, "-i", s, p)

		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr

		if err := cmd.Run(); err != nil {
			fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
			return err
		}
	}

	return nil
}

func main() {
	s := flag.String("s", "", "Source directory.")
	o := flag.String("o", "", "Output directory.")
	c := make(chan string)

	flag.Parse()

	var wg sync.WaitGroup

	go func() {
		if err := Scan(*s, c, ".mp4 | .mov"); err != nil {
			log.Fatalf("error scanning directory: %v", err)
		}

		wg.Wait()
		close(c)
	}()

	for x := range c {
		wg.Add(1)
		go func(s string) {
			fmt.Printf("\nExtracting audio from %s...", s)
			if err := Extract(s, "ffmpeg", *o); err != nil {
				log.Fatalf("error extracting audio: %v", err)
			}

			fmt.Printf("\n%s Done!", s)
			wg.Done()
		}(x)
	}

	fmt.Print("\n")
}
