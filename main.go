package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type Episode struct {
	Filename string
	Number   int
	Title    string
}

func main() {
	flag.Parse()
	wd := workingDir(flag.Args())

	episodes := []Episode{}
	fmt.Println("Please enter the episode number of each file.")
	for _, file := range listFiles(wd, recursive) {
		fmt.Println(file)
		str := input("> ")
		if str == "--" {
			break
		}
		if str != "-" {
			number, err := strconv.Atoi(str)
			if err != nil {
				log.Fatal(err)
			}
			ep := Episode{Filename: filepath.Join(wd, file), Number: number}
			episodes = append(episodes, ep)
		}
	}

	sort.Slice(episodes, func(i, j int) bool { return episodes[i].Number < episodes[j].Number })
	fmt.Println("Please enter the episode title for each file.")
	for i, ep := range episodes {
		str := input(fmt.Sprintf("Episode %d: ", ep.Number))
		if !original {
			str = correctFilename(str)
		}
		episodes[i].Title = str
	}

	for _, ep := range episodes {
		dest := fmt.Sprintf(pattern, ep.Number, ep.Title) + filepath.Ext(ep.Filename)
		dest, err := filepath.Abs(dest)
		if err != nil {
			log.Fatal(err)
		}
		err = os.Rename(ep.Filename, dest)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func workingDir(relPaths []string) string {
	if len(relPaths) >= 1 {
		dir, err := filepath.Abs(relPaths[0])
		if err != nil {
			log.Fatal(err)
		}
		fileInfo, err := os.Stat(dir)
		if err != nil {
			log.Fatal(err)
		}
		if fileInfo.IsDir() {
			return dir
		} else {
			return filepath.Dir(dir)
		}
	} else {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		return dir
	}
}

func listFiles(dir string, recursive bool) []string {
	matches, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		log.Fatal(err)
	}

	files := []string{}
	dirs := []string{}
	for _, file := range matches {
		fileInfo, err := os.Stat(file)
		if err != nil {
			log.Fatal(err)
		}
		if fileInfo.Mode().IsRegular() {
			files = append(files, file)
		} else if recursive && fileInfo.IsDir() {
			dirs = append(dirs, file)
		}
	}
	for _, d := range dirs {
		for _, file := range listFiles(d, true) {
			files = append(files, filepath.Join(d, file))
		}
	}

	for i := range files {
		filename, err := filepath.Rel(dir, files[i])
		if err != nil {
			log.Fatal(err)
		}
		files[i] = filename
	}
	return files
}

func correctFilename(filename string) string {
	forbidens := []string{"<", ">", ":", "\"", "/", "\\", "|", "?", "*"}
	for _, forbiden := range forbidens {
		filename = strings.ReplaceAll(filename, forbiden, "_")
	}
	return filename
}
