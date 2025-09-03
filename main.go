package main

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	Editor     string
	NotesPath  string
	NotePeriod string
}

var conf = Config{
	Editor:     "/usr/bin/nvim",
	NotesPath:  "~/.notes/",
	NotePeriod: "day", // day or week
}

func newNoteFile(filepath string) error {
	return nil
}

func createFilepath(conf Config) (string, error) {
	now := time.Now()
	var filename string

	switch conf.NotePeriod {
	case "day":
		filename = now.Format("day_2006-01-02")
	case "week":
		// Find the most recent Monday
		weekday := now.Weekday()
		mondayOffset := (weekday + 6) % 7 // Days to subtract to get to Monday
		monday := now.AddDate(0, 0, -int(mondayOffset))
		filename = monday.Format("week_2006-01-02")
	default:
		return "", os.ErrInvalid
	}

	fullPath := filepath.Join(conf.NotesPath, filename)

	// Check if file already exists
	if _, err := os.Stat(fullPath); err != nil {
		return fullPath, os.ErrExist
	}

	return fullPath, nil
}

func main() {
	filepath, err := createFilepath(conf)
	if err != nil {
		log.Panic(err)
	}
	log.Println(filepath)
	// check if there is a file for day or week in the path
	// if not, create it
	// write current timestamp down
	// open editor with file

	// more command line options
	// -c     just dump the file instead of editing
	// -b x   go back x files and edit/output this
}
