package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/larry868/linkerpod/pkg/yamlpod"
	"github.com/larry868/loadfavicon/v2"
	"github.com/larry868/verbose"
)

const (
	PATH_CACHEICONS string = "assets/.cacheicons/"
)

// linkerpod command line tool used to load link's favicons and to set it in the static website cache.
//
// # parameters:
//   - pod: full path of the yaml file selected. "linkerpod.yaml" by default.
//   - resetcache: reset favicon cache of the pod. false by default.
//   - loadfavicons: load and set pod's favicons in cache. false by default.
//   - overwrite: overwrite favicons already in cache. false by default.
//   - verbose: verbose output.
//   - debug: generate more output for debugging.
func main() {
	var pod string
	var resetcache, loadicons, overwrite bool

	flag.StringVar(&pod, "pod", "linkerpod.yaml", "full path of the yaml file selected")
	flag.BoolVar(&resetcache, "resetcache", false, "reset favicon cache of the pod")
	flag.BoolVar(&loadicons, "loadfavicons", false, "load and set pod's favicons in cache")
	flag.BoolVar(&overwrite, "overwrite", false, "overwrite favicons already in cache")
	flag.BoolVar(&verbose.IsOn, "verbose", false, "verbose output")
	flag.BoolVar(&verbose.IsDebugging, "debug", false, "generate more output for debugging")
	flag.Parse()

	// get the yaml file for this pod
	if pod == "" {
		fmt.Println("the pod parameter is required")
		os.Exit(1)
	}

	var yamlpath, yamlname string
	pext := path.Ext(pod)
	if pext != "" && pext != ".yaml" {
		fmt.Println("the pod parameter must refer to a yaml file")
		os.Exit(-1)
	} else if pext == ".yaml" {
		yamlname = path.Base(pod)[:len(path.Base(pod))-5]
		yamlpath = pod
	} else {
		yamlname = path.Base(pod)
		yamlpath = pod + ".yaml"
	}

	_, err := os.Stat(yamlpath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cachepath := path.Join(filepath.Dir(yamlpath), PATH_CACHEICONS)
	if resetcache {
		os.RemoveAll(cachepath)
		fmt.Printf("favicon cache %q has been reset\n", cachepath)
	}

	file, err := os.OpenFile(yamlpath, os.O_RDWR, os.ModeExclusive|os.ModePerm) // 0x755
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	ys, err := yamlpod.Unmarshal(yamlpath, file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// update cache
	if loadicons || resetcache {
		verbose.Printf(verbose.INFO, "processing pod %q: %v links\n", yamlpath, len(ys.SingleLinks))
		nko := 0
		nok := 0
		client := &http.Client{Timeout: time.Second * 5}
		for x, ylnk := range ys.SingleLinks {
			ylnk.Link = strings.ToLower(strings.Trim(ylnk.Link, " "))

			isincache := strings.HasPrefix(ylnk.Icon, PATH_CACHEICONS)

			fup := false
			if resetcache {
				if isincache {
					ylnk.Icon = ""
					fup = true
					isincache = false
				}
			}

			isblank := ylnk.Icon == ""

			if loadicons && ylnk.Link != "" {
				if isblank || (overwrite && isincache) {
					iconfname, err := loadfavicon.DownloadOne(client, ylnk.Link, cachepath, overwrite)
					if err != nil {
						verbose.Println(verbose.WARNING, err.Error())
						nko++
					} else if iconfname != "" {
						ylnk.Icon = path.Join(PATH_CACHEICONS, iconfname)
						fup = true
					}
				}
			}

			if fup {
				nok += 1
				ys.SingleLinks[x] = ylnk
			}
		}

		if nok > 0 {
			// TODO: backup previous file
			file.Seek(0, 0)
			file.Truncate(0)
			file.WriteString("# Linkerpod setup file\n")
			file.WriteString(fmt.Sprintf("# %s updated on %s\n", yamlname+".yaml", time.Now().Format("2006-01-02 15:04:05")))
			err := yamlpod.Marshal(file, ys)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		fmt.Printf("favicon cache: %v link update, %v link fail\n", nok, nko)
		os.Exit(0)
	} else {
		fmt.Printf("pod %q: has %v minipods and %v links\n", yamlpath, len(ys.MiniPods), len(ys.SingleLinks))
	}

}
