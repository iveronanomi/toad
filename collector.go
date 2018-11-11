package todo

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

const defaultChunkSize = 1024 * 1024

func New(sources, excludes, extensions []string, chunkSize int64) *walker {
	if chunkSize < 1 {
		chunkSize = defaultChunkSize
	}
	return &walker{
		sources:    sources,
		extensions: extensions,
		excludes:   excludes,
		chunkSize:  chunkSize,
	}
}

type walker struct {
	sources    []string
	excludes   []string
	extensions []string
	chunkSize  int64

	found uint64
}

func (w *walker) Collect(chunks chan []string, filterFunc ...FilterFunc) error {

	var size int64
	var chunk []string
	var chunksEncounter int

	for _, src := range w.sources {
		err := filepath.Walk(src,
			func(path string, info os.FileInfo, err error) error {
				if info == nil {
					return nil
				}
				if info.IsDir() {
					return nil
				}
				if info.Size() < 1 {
					return nil
				}

				for _, exclude := range w.excludes {
					if len(path) > len(exclude) {
						continue
					}
					if strings.HasPrefix(path, exclude) {
						return nil
					}
				}

				var n int
				if len(w.extensions) > 0 {
					ext := filepath.Ext(path)
					if len(ext) < 2 {
						return nil
					}
					ext = ext[1:]
					for _, extension := range w.extensions {
						if ext == extension {
							n++
							break
						}
					}
					if n < 1 {
						return nil
					}
				}
				for i := range filterFunc {
					if filterFunc[i](path) {
						return nil
					}
				}
				if size+info.Size() > w.chunkSize {
					chunks <- chunk
					chunk = []string{path}
					size = info.Size()
					chunksEncounter++
					return nil
				}
				w.found++
				chunk = append(chunk, path)
				size += info.Size()
				return err
			})
		if err != nil {
			log.Println(err)
		}
	}

	chunks <- chunk
	chunksEncounter++
	chunks <- []string{}
	close(chunks)

	log.Printf("chunks %v", chunksEncounter)
	return nil
}

type FilterFunc func(path string) bool
