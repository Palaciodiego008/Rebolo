package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

func startDevServer() {
	// Start Bun in watch mode for assets
	go startBunWatcher()
	
	// Start Go server with hot reload
	startGoServer()
}

func startBunWatcher() {
	if _, err := os.Stat("package.json"); os.IsNotExist(err) {
		fmt.Println("No package.json found, skipping Bun watcher")
		return
	}
	
	fmt.Println("🟡 Starting Bun asset watcher...")
	cmd := exec.Command("bun", "run", "dev")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Start(); err != nil {
		log.Printf("Failed to start Bun: %v", err)
	}
}

func startGoServer() {
	fmt.Println("🔥 Starting Go server with hot reload...")
	
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	
	// Watch Go files
	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if info.IsDir() && shouldSkipDir(path) {
			return filepath.SkipDir
		}
		
		if strings.HasSuffix(path, ".go") {
			watcher.Add(filepath.Dir(path))
		}
		
		return nil
	})
	
	var cmd *exec.Cmd
	restartServer := func() {
		if cmd != nil && cmd.Process != nil {
			cmd.Process.Kill()
			cmd.Wait()
		}
		
		fmt.Println("🔄 Restarting server...")
		cmd = exec.Command("go", "run", "main.go")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Start()
	}
	
	// Initial start
	restartServer()
	
	// Watch for changes
	debounce := time.NewTimer(100 * time.Millisecond)
	debounce.Stop()
	
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			
			if event.Op&fsnotify.Write == fsnotify.Write && strings.HasSuffix(event.Name, ".go") {
				debounce.Reset(100 * time.Millisecond)
			}
			
		case <-debounce.C:
			restartServer()
			
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("Watcher error:", err)
		}
	}
}

func shouldSkipDir(path string) bool {
	skipDirs := []string{"node_modules", ".git", "vendor", "public"}
	for _, skip := range skipDirs {
		if strings.Contains(path, skip) {
			return true
		}
	}
	return false
}
