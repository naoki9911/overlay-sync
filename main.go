package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

var threadLock = make(chan bool, 20)
var threadWg sync.WaitGroup

func main() {
	if len(os.Args) < 3 {
		panic("./overlay-sync mounted upper")
	}

	mountedPath := os.Args[1]
	upperPath := os.Args[2]
	err := syncDir(mountedPath, upperPath, "")
	if err != nil {
		panic(err)
	}

	threadWg.Wait()
}

func syncDir(mountedPath, upperPath, relPath string) error {
	mount := filepath.Join(mountedPath, relPath)
	files, err := os.ReadDir(mount)
	if err != nil {
		return err
	}

	for _, file := range files {
		upper := filepath.Join(upperPath, relPath, file.Name())
		_, err = os.Stat(upper)
		if file.IsDir() {
			if err != nil {
				if !os.IsNotExist(err) {
					return err
				}
				threadLock <- true
				threadWg.Add(1)
				go func(path string) {
					err = syncNotExistingDir(path)
					if err != nil {
						panic(err)
					}
					<-threadLock
					threadWg.Done()
				}(filepath.Join(mount, relPath, file.Name()))
			} else {
				threadLock <- true
				threadWg.Add(1)
				go func(relPath string) {
					err = syncDir(mountedPath, upperPath, relPath)
					if err != nil {
						panic(err)
					}
					<-threadLock
					threadWg.Done()
				}(filepath.Join(relPath, file.Name()))
			}
		} else {
			fmt.Printf("%s is checked\n", upper)
			if err == nil {
				continue
			}
			if !os.IsNotExist(err) {
				return err
			}

			syncPath := filepath.Join(mount, file.Name())
			f, err := os.OpenFile(syncPath, os.O_WRONLY|os.O_APPEND, os.ModePerm)
			if err != nil {
				fmt.Printf("failed to open file %v\n", err)
				// ignore error
				continue
			}
			f.Close()
			fmt.Printf("%s is synced\n", syncPath)
		}
	}

	return nil
}

func syncNotExistingDir(path string) error {
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	for _, file := range files {
		syncPath := filepath.Join(path, file.Name())
		if file.IsDir() {
			threadLock <- true
			threadWg.Add(1)
			go func(path string) {
				err = syncNotExistingDir(path)
				if err != nil {
					panic(err)
				}
				<-threadLock
				threadWg.Done()
			}(syncPath)
		}

		f, err := os.OpenFile(syncPath, os.O_WRONLY|os.O_APPEND, os.ModePerm)
		if err != nil {
			return err
		}
		f.Close()
		fmt.Printf("%s is synced\n", syncPath)
	}

	return nil
}
