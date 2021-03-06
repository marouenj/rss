package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/marouenj/rss/agent"
)

var (
	config string = "config"
	data   string = "data"
	in     string = "channels"
	out    string = "items"
)

func main() {
	// parse args
	baseDir := flag.String("base_dir", "./", "")
	flag.Parse()

	// check baseDir exists
	if _, err := os.Stat(*baseDir); err != nil {
		fmt.Printf("[ERR] Base dir not exists: %v\n", err)
		os.Exit(1)
	}

	// check baseDir is a dir
	info, _ := os.Stat(*baseDir)
	if !info.IsDir() {
		fmt.Printf("[ERR] Base dir not a dir\n")
		os.Exit(1)
	}

	dataDir := filepath.Join(*baseDir, data)

	// check dataDir exists
	if _, err := os.Stat(dataDir); err != nil {
		fmt.Printf("[ERR] Data dir not exists: %v\n", err)
		os.Exit(1)
	}

	// check dataDir is a dir
	info, _ = os.Stat(dataDir)
	if !info.IsDir() {
		fmt.Printf("[ERR] Data dir not a dir\n")
		os.Exit(1)
	}

	inDir := filepath.Join(dataDir, in)

	// check inDir exists
	if _, err := os.Stat(inDir); err != nil {
		fmt.Printf("[ERR] Input dir not exists: %v\n", err)
		os.Exit(1)
	}

	// check inDir is a dir
	info, _ = os.Stat(inDir)
	if !info.IsDir() {
		fmt.Printf("[ERR] Input dir not a dir\n")
		os.Exit(1)
	}

	outDir := filepath.Join(dataDir, out)

	// check outDir exists
	if _, err := os.Stat(outDir); err != nil {
		// create if not exists
		err = os.Mkdir(outDir, os.ModeDir|os.ModePerm)
		if err != nil {
			fmt.Printf("[ERR] Unable to create dir '%s': %v\n", outDir, err)
			os.Exit(1)
		}
	}

	// check outDir is a dir
	info, _ = os.Stat(outDir)
	if !info.IsDir() {
		fmt.Printf("[ERR] Output dir not a dir\n")
		os.Exit(1)
	}

	// create loader
	loader, err := agent.NewLoader()
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	// load
	err = loader.Load(inDir)
	if err != nil {
		fmt.Printf("[ERR] Unable to load channels: %v\n", err)
		os.Exit(1)
	}

	// create crawler
	crawler, err := agent.NewCrawler()
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	// crawl
	err = crawler.Crawl(loader)
	if err != nil {
		fmt.Printf("[ERR] Unable to download items: %v\n", err)
		os.Exit(1)
	}

	// create marshaller
	marshaller, err := agent.NewMarshaller(outDir)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	// rearrange items
	err = marshaller.ReArrange(crawler.Rss.Channels)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	// persist
	err = marshaller.Save()
	if err != nil {
		fmt.Printf("[ERR] Unable to merge and persist new items: %v\n", err)
		os.Exit(1)
	}
}
