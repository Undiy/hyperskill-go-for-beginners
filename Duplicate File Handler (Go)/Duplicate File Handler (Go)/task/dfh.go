package main

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type PathsWithSize struct {
	paths []string
	size  int64
}

type DuplicatePathsWithSize struct {
	paths []DuplicatePaths
	size  int64
}

type DuplicatePaths struct {
	paths []DuplicatePath
	hash  string
}

type DuplicatePath struct {
	path  string
	index int
}

func addPath(files *[]PathsWithSize, path string, size int64) {
	for i, pws := range *files {
		if pws.size == size {
			(*files)[i].paths = append(pws.paths, path)
			return
		}
	}
	*files = append(*files, PathsWithSize{[]string{path}, size})
}

func calcHash(path string) string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	md5Hash := md5.New()
	if _, err = io.Copy(md5Hash, file); err != nil {
		log.Fatal(err)
	}

	s := fmt.Sprintf("%x", md5Hash.Sum(nil))
	return s
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Directory is not specified")
		return
	}

	root := os.Args[1]

	var ext string
	fmt.Println("\nEnter file format:")
	fmt.Scanln(&ext)

	fmt.Print(`
Size sorting options:
1. Descending
2. Ascending
`)

	var isAsc = false
	for {
		fmt.Println("\nEnter a sorting option:")
		var option int
		fmt.Scan(&option)
		if option == 1 || option == 2 {
			isAsc = option == 2
			break
		} else {
			fmt.Println("\nWrong option")
		}
	}

	var files = make([]PathsWithSize, 0)
	filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
			return err
		}
		if !info.IsDir() && (ext == "" || filepath.Ext(path) == "."+ext) {
			addPath(&files, path, info.Size())
		}

		return nil
	})

	sort.Slice(files, func(i, j int) bool {
		if isAsc {
			return files[i].size < files[j].size
		} else {
			return files[i].size > files[j].size
		}
	})

	fmt.Println("")
	for _, pws := range files {
		fmt.Printf("%d bytes\n%s\n\n", pws.size, strings.Join(pws.paths, "\n"))
	}

	fmt.Println("\nCheck for duplicates?")
	var check string
	for check != "yes" {
		fmt.Scanln(&check)
		if check == "no" {
			return
		} else if check != "yes" {
			fmt.Println("Wrong option")
		}
	}

	duplicatesWithSize := make([]DuplicatePathsWithSize, 0)
	index := 1

	for _, pws := range files {
		pathsByHash := make(map[string][]string, 0)
		for _, p := range pws.paths {
			hash := calcHash(p)
			duplicates, ok := pathsByHash[hash]
			if ok {
				pathsByHash[hash] = append(duplicates, p)
			} else {
				pathsByHash[hash] = []string{p}
			}
		}

		dws := DuplicatePathsWithSize{make([]DuplicatePaths, 0), pws.size}
		for hash, ds := range pathsByHash {
			if len(ds) < 2 {
				continue
			}
			dps := DuplicatePaths{make([]DuplicatePath, 0), hash}
			for _, path := range ds {
				dps.paths = append(dps.paths, DuplicatePath{path, index})
				index++
			}
			dws.paths = append(dws.paths, dps)
		}
		if len(dws.paths) > 0 {
			duplicatesWithSize = append(duplicatesWithSize, dws)
		}
	}

	for _, dws := range duplicatesWithSize {
		fmt.Printf("\n%d bytes\n", dws.size)
		for _, dps := range dws.paths {
			fmt.Printf("Hash: %s\n", dps.hash)
			for _, dp := range dps.paths {
				fmt.Printf("%d. %s\n", dp.index, dp.path)
			}
		}
	}

	fmt.Println("\nDelete files?")
	var delete string
	for delete != "yes" {
		fmt.Scanln(&delete)
		if delete == "no" {
			return
		} else if delete != "yes" {
			fmt.Println("\nWrong option")
		}
	}

	reader := bufio.NewReader(os.Stdin)
	var filesToDelete map[int]bool
	for filesToDelete == nil {
		fmt.Println("\nEnter file numbers to delete:")
		filesToDeleteStr, _ := reader.ReadString('\n')
		filesToDeleteStr = strings.TrimSpace(filesToDeleteStr)
		for _, s := range strings.Split(filesToDeleteStr, " ") {
			n, err := strconv.Atoi(s)
			if err != nil {
				fmt.Println("\nWrong format")
				filesToDelete = nil
				break
			}
			if filesToDelete == nil {
				filesToDelete = map[int]bool{n: true}
			} else {
				filesToDelete[n] = true
			}
		}
		fmt.Println("Done with ", filesToDeleteStr)
	}

	fmt.Println(filesToDelete)

	var freed int64 = 0
	for _, dws := range duplicatesWithSize {
		for _, dps := range dws.paths {
			for _, dp := range dps.paths {
				if filesToDelete[dp.index] {
					err := os.Remove(dp.path)
					if err != nil {
						log.Fatal(err)
					}
					freed += dws.size
				}
			}
		}
	}

	fmt.Printf("\nTotal freed up space: %d bytes", freed)
}
