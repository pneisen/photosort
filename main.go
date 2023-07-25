package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
)

// getDataTime uses exif to get the date and time a photo was taken.
func getDateTime(path string) (time.Time, error) {
	f, err := os.Open(path)
	if err != nil {
		return time.Time{}, err
	}
	defer f.Close()

	exif.RegisterParsers(mknote.All...)

	x, err := exif.Decode(f)
	if err != nil {
		return time.Time{}, err
	}

	tm, err := x.DateTime()
	if err != nil {
		return time.Time{}, err
	}

	return tm, nil
}

// dirExists only returns true if the directory can be stat'd and IS a directory.
func dirExists(path string) bool {
	fi, err := os.Stat(path)
	if err == nil {
		if fi != nil {
			return fi.IsDir()
		}
		return false
	}
	return false
}

// fileExists only returns true if the file can be stat'd and is not a directory.
func fileExists(path string) bool {
	fi, err := os.Stat(path)
	if err == nil {
		if fi != nil {
			return !fi.IsDir()
		}
		return false
	}
	return false
}

// cp copies a file with limited error checks.
func cp(fromPath, toPath string) error {
	from, err := os.Open(fromPath)
	if err != nil {
		return err
	}

	defer from.Close()

	to, err := os.Create(toPath)
	if err != nil {
		return err
	}

	defer to.Close()
	_, err = io.Copy(to, from)
	if err != nil {
		return err
	}

	return to.Close()
}

func main() {
	var err error
	var fromPath string
	var toPath string

	if len(os.Args) < 3 {
		log.Fatal("Usage: photosort <from dir> <to dir>")
	} else {
		fromPath = os.Args[1]
		toPath = os.Args[2]
	}

	// make sure paths are both valid
	_, err = os.Stat(fromPath)
	if err != nil {
		log.Fatal("Error: Specified directory not valid")
	}
	_, err = os.Stat(toPath)
	if err != nil {
		log.Fatal("Error: Specified directory not valid")
	}

	unsortedDirName := filepath.Join(toPath, "unsorted")

	// Walk the path tree, figure out the date pics were taken, and move files into YYYY-Month named directories.
	err = filepath.Walk(fromPath, func(path string, info os.FileInfo, err error) error {

		// We only care about the files.
		if !info.IsDir() {
			tm, err := getDateTime(path)
			if err != nil {
				// Something went wrong and we don't have a date for this pic, throw it in unsorted.
				if !dirExists(unsortedDirName) {
					os.Mkdir(unsortedDirName, os.FileMode(0777))
				}

				if fileExists(filepath.Join(unsortedDirName, filepath.Base(path))) {
					println("Error: File exists. Not Copied to " + unsortedDirName + ": " + path)
				} else {
					err := cp(path, filepath.Join(unsortedDirName, filepath.Base(path)))
					if err != nil {
						log.Fatal("Error: A copy failed: " + path)
					}
				}

			} else {
				dateString := strconv.Itoa(tm.Year()) + "-" + strconv.Itoa(int(tm.Month()))
				dateDirName := filepath.Join(toPath, dateString)

				if !dirExists(dateDirName) {
					os.Mkdir(dateDirName, os.FileMode(0777))
				}

				if fileExists(filepath.Join(dateDirName, filepath.Base(path))) {
					println("Error: File exists. Not Copied to " + dateString + ": " + path)
				} else {
					err := cp(path, filepath.Join(dateDirName, filepath.Base(path)))
					if err != nil {
						log.Fatal("Error: A copy failed: " + path)
					}
				}
			}
		}

		return nil
	})

}
