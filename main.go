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

func main() {
	s := flag.String("s", "", "Source directory. (Ex. -s=~/Videos)")
	o := flag.String("o", "", "Output directory. (Ex. -o=~/Music)")
	x := flag.String("x", ".mp4 | .mov", "Video files extensions.")

	c := make(chan string, 4)

	flag.Parse()

	if *s == "" && *o == "" {
		fmt.Println("Source and output directory are required. Type `mp3xtrak -h` for help.")
		return
	}

	var wg sync.WaitGroup

	go func() {
		if err := scan(&wg, *s, c, *x); err != nil {
			log.Fatalf("error scanning directory: %v", err)
		}

		wg.Wait()
		close(c)
	}()

	for v := range c {
		go func(s string) {
			fmt.Printf("\nextracting audio from %s...", s)
			if err := extract(&wg, s, "ffmpeg", *o); err != nil {
				fmt.Printf("\nerror extracting audio: %v", err)
			} else {
				fmt.Printf("\n%s Done!", s)
			}
		}(v)
	}

	fmt.Print("\n")
}

func scan(wg *sync.WaitGroup, dir string, c chan string, exts string) error {
	err := filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.ContainsAny(info.Name(), exts) {
			wg.Add(1)
			c <- path.Join(dir, info.Name())
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func extract(wg *sync.WaitGroup, s string, command, dir string) error {
	defer wg.Done()
	if s != "" {
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

		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("%v: %s", err, stderr.String())
		}
	}

	return nil
}

func escape(s string) string {
	chars := " ()!$&'*,;<=>?[]^`{}|~"
	for _, c := range []byte(chars) {
		v := string(c)
		s = strings.Replace(s, v, "\\"+v, -1)
	}
	return s
}
