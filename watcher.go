package main

import (
	"log"
	"strings"
	"time"

	"github.com/howeyc/fsnotify"
)

const (
	// Because of multiple events following one another in succession, lets delay the regeneration for some time
	DELAY = 10 * time.Second
)

// Watch a specific directory
func watchDir(dir string, watcher *fsnotify.Watcher) {
	if err := watcher.Watch(dir); err != nil {
		watcher.Close()
		log.Fatal("[ERROR]", err)
	}
}

// Starts the watching thread
func startWatcher() *fsnotify.Watcher {
	watcher, err := fsnotify.Watcher()
	if err != nil {
		log.Fatal("[ERROR]", err)
	}
	go watch(watcher)

	// Watch posts directory
	watchDir(PostsDir, watcher)
	// Watch templates directory
	watchDir(TemplatesDir, watcher)
	return watcher
}

// Based on the event ping for an engine re run
func watch(watcher *fsnotify.Watcher) {
	var ping <-chan time.Time
	for {
		select {
		case event := <-watcher.Event:
			// Only ping for a re run of the engine if the changes affect markdown files or templates
			ext := event.Name
			if strings.HasPrefix(ext, PostsDir) && ext == ".md" {
				ping = time.After(DELAY)
			} else if strings.HasPrefix(ext, TemplatesDir) && ext == ".html" {
				ping = time.After(DELAY)
			}
		case err := <-watcher.Error:
			log.Println("[INFO] Watcher error", err)
		case <-ping:
			if err := Run(); err != nil {
				log.Println("[INFO] Error running engine: ", err)
			} else {
				log.Println("[SUCCESS] Engine run")
			}
		}
	}
}
