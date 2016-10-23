package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh/terminal"
)

// This is heavily based on gopl.io/ch8/du4

type fileInfo struct {
	path string // In contrast to os.FileInfo, this contains the whole path
	size int64
}

// Sorting entries
type bySize []fileInfo

// sort Interface
func (s bySize) Len() int           { return len(s) }
func (s bySize) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s bySize) Less(i, j int) bool { return s[i].size < s[j].size }

var Program, Version string

var (
	flagN  = flag.Int("n", 10, "Print this `number` of largest files")
	flagO  = flag.Bool("o", false, "Print file names only (e.g., for xargs scripting)")
	flagC  = flag.Bool("c", false, "Files from current directory only (no recursion)")
	flagV  = flag.Bool("version", false, "Print version and exit")
	flagCC = flag.Int("concurrent", 20, "Start this `number` of concurrent tree walks (values <= 0 get set to 1)")
	flagU  = UnitFlag("u", Units["GiB"], "Print sizes in specified `unit` ("+allUnits()+")")
	flagI  = IgnoreFlag("i", Ignores{}, "Ignore files/directories matching `regex` (can be used multiple times)")
)

var done = make(chan struct{})

func main() {
	Program = path.Base(os.Args[0])
	if Version == "" {
		Version = "Unknown Version"
	}

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s (%s):\n", Program, Version)
		flag.PrintDefaults()
	}

	flag.Parse()

	if *flagV {
		fmt.Println(Program, Version)
		return
	}

	if *flagCC <= 0 {
		*flagCC = 1
	}

	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}

	// Cancel traversal when input is detected.
	fmt.Fprintf(os.Stderr, "Press <return> to stop processing\n")
	go func() {
		os.Stdin.Read(make([]byte, 1)) // Read a single byte.
		close(done)
	}()

	// Traverse each root of the file tree in parallel.
	fileInfos := make(chan fileInfo)
	var n sync.WaitGroup
	for _, root := range roots {
		n.Add(1)
		go walkDir(root, &n, fileInfos)
	}
	go func() {
		n.Wait()
		close(fileInfos)
	}()

	// Print the results periodically.
	tick := time.Tick(500 * time.Millisecond)
	var nFiles, nBytes int64

	largest := make([]fileInfo, *flagN)
	for i := 0; i < *flagN; i++ {
		largest[i].size = -1
	}

loop:
	for {
		select {
		case <-done:
			// Drain fileInfos to allow existing goroutines to finish.
			for range fileInfos {
				// Do nothing.
			}
			return
		case info, ok := <-fileInfos:
			if !ok {
				break loop // fileInfos was closed (all workers finished).
			}
			nFiles++
			nBytes += info.size
			if info.size > largest[0].size && !ignore(info.path) {
				largest[0] = info
				sort.Sort(bySize(largest))
			}
		case <-tick:
			fmt.Fprintf(os.Stderr, "Processed: %d files  %.1f %s\n", nFiles, float64(nBytes)/float64(*flagU), *flagU)
		}
	}
	printDiskUsage(&largest, nBytes)
}

func cancelled() bool {
	select {
	case <-done:
		return true
	default:
		return false
	}
}

func walkDir(dir string, n *sync.WaitGroup, fileInfos chan<- fileInfo) {
	defer n.Done()
	if cancelled() {
		return
	}
	for _, entry := range dirents(dir) {
		fullPath := filepath.Join(dir, entry.Name())
		if entry.IsDir() {
			if !*flagC && !ignore(fullPath) {
				n.Add(1)
				go walkDir(fullPath, n, fileInfos) // walk sub directory
			}
		} else {
			fileInfos <- fileInfo{path: fullPath, size: entry.Size()}
		}
	}
}

var sema = make(chan struct{}, *flagCC) // concurrency-limiting counting semaphore

func dirents(dir string) []os.FileInfo {
	select {
	case sema <- struct{}{}: // acquire token
	case <-done:
		return nil // cancelled
	}
	defer func() { <-sema }() // release token

	f, err := os.Open(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", Program, err)
		return nil
	}
	defer f.Close()

	entries, err := f.Readdir(0) // 0 => no limit; read all entries
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", Program, err)
		// Don't return: Readdir may return partial results.
	}

	return entries
}

// we call ignore for
// - a directory in walkDir (cheaper than calling for every file, but allows cutting whole directory branches)
// - when we insert a file if it is big enough (usually also cheap, because does not happen that often)
func ignore(path string) bool {
	for _, r := range *flagI {
		if ignore := r.MatchString(path); ignore {
			return true
		}
	}
	return false
}

func printDiskUsage(largest *[]fileInfo, nBytes int64) {
	width, _, _ := terminal.GetSize(0)
	w := width - 2 // "[]"
	if w == -2 {   // could not determine terminal width, eg Docker
		w = 78
	}

	max := (*largest)[len(*largest)-1].size

	var sum int64
	var entries int

	for i := len(*largest) - 1; i >= 0; i-- {
		f := (*largest)[i]
		if f.size != -1 {
			sum += f.size
			entries++
			if *flagO {
				fmt.Println(f.path)
				continue
			}

			fmt.Printf("%.1f %s: %s\n", float64(f.size)/float64(*flagU), *flagU, f.path)
			var h int
			if max != 0 {
				h = int(int64(w) * f.size / max)
			} else {
				h = w
			}
			fmt.Printf("[%s]\n\n", strings.Repeat("#", h))
		}
	}
	if !*flagO {
		var p string
		if entries > 1 {
			p = "s"
		}
		fmt.Printf("%d largest file%s: %.1f %s / %.1f %s total\n", entries, p,
			float64(sum)/float64(*flagU), *flagU,
			float64(nBytes)/float64(*flagU), *flagU)
	}
}
