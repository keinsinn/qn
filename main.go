package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Config struct {
	Editor     string
	NotesPath  string
	NotePeriod string
}

// FIXME: This needs a config file
var conf = Config{
	Editor:     "/usr/bin/nvim",
	NotesPath:  "~/.notes/",
	NotePeriod: "day", // day or week
}

func replaceHome(input string) string {
	return strings.ReplaceAll(input, "~", os.Getenv("HOME"))
}

func newNoteFile(fullPath string) error {
	now := time.Now()
	content := []byte("# Notizen ab dem " + now.Format("02.01.2006") + "\n")
	err := os.WriteFile(fullPath, content, 0o644)
	return err
}

func insertTimestamp(fullPath string) error {
	now := time.Now()
	content := "\n## " + now.Format("02.01.2006 - 15:04") + "\n\n"

	f, err := os.OpenFile(fullPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o600)
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err = f.WriteString(content); err != nil {
		return err
	}
	return nil
}

func fileExists(fullPath string) bool {
	if _, err := os.Stat(fullPath); err == nil {
		return true
	}
	return false
}

func createFilepath(conf Config) (string, error) {
	now := time.Now()
	var filename string

	switch conf.NotePeriod {
	case "day":
		filename = now.Format("day_2006-01-02.md")
	case "week":
		// Find the most recent Monday
		weekday := now.Weekday()
		mondayOffset := (weekday + 6) % 7 // Days to subtract to get to Monday
		monday := now.AddDate(0, 0, -int(mondayOffset))
		filename = monday.Format("week_2006-01-02.md")
	default:
		return "", os.ErrInvalid
	}

	fullPath := filepath.Join(replaceHome(conf.NotesPath), filename)

	return fullPath, nil
}

func runEditor(filepath string) error {
	// FIXME: arguments for neovim, needs more generalization
	cmd := exec.Command(conf.Editor, "+startinsert", "+", filepath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	return err
}

func main() {
	filepath, err := createFilepath(conf)
	if err != nil {
		log.Panic(err)
	}
	if !fileExists(filepath) {
		newNoteFile(filepath)
	}
	insertTimestamp(filepath)
	err = runEditor(filepath)
	if err != nil {
		log.Panic(err)
	}

	// more command line options
	// -c     just dump the file instead of editing
	// -b x   go back x files and edit/output this
}
