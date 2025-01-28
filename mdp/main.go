package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

type content struct {
	Title string
	Body  template.HTML
}

const (
	defaultTemplate = `<!DOCTYPE html>
	<html>
	<head>
	<meta http-equiv="content-type" content="text/html; charset=utf-8">
	<title>{{ .Title }}</title>
	</head>
	<body>
	{{ .Body }}
	</body>
	</html>
	`
)

func main() {
	times := 0
	for {
		times++
		counter := PackItems(0)
		if counter != 2000 {
			log.Fatalf("it should be 2000 but found %d on execution %d", counter, times)
		}
	}
}

// func mainPackItems() {
// 	fmt.Println("Total items packed: ", PackItems(0))
// }

func PackItems(totalItems int32) int32 {
	const workers = 2
	const itemsPerWorker = 1000

	var wg sync.WaitGroup
	itemsPacked := 0
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for j := 0; j < itemsPerWorker; j++ {
				// simulate packing an item
				//atomic structures are awesome when we need to sync access to a single operation
				//for group of operations its better to use sync.Mutex
				atomic.AddInt32(&totalItems, int32(itemsPacked))
				//update total items packed without proper sync
				// totalItems = itemsPacked
			}
		}(i)
	}
	wg.Wait()
	return totalItems
}

func main2() {
	filename := flag.String("file", "", "Markdown file to preview")
	skipPreview := flag.Bool("s", false, "Skip preview")
	tFname := flag.String("t", "", "Alternate template name")
	flag.Parse()

	if *filename == "" {
		flag.Usage()
		os.Exit(1)
	}

	if err := run(*filename, *tFname, os.Stdout, *skipPreview); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(filename string, tFname string, out io.Writer, skipPreview bool) error {
	input, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	htmlData, err := parseContent(input, tFname)
	if err != nil {
		return err
	}

	temp, err := os.CreateTemp("", "mdp_*.html")
	if err != nil {
		return err
	}
	if err := temp.Close(); err != nil {
		return err
	}

	outName := temp.Name()
	fmt.Fprintln(out, outName)

	if err := saveHTML(outName, htmlData); err != nil {
		return err
	}
	if skipPreview {
		return nil
	}
	defer os.Remove(outName)
	return preview(outName)
}

func checkConfigurableTemplate(t template.Template, tFname string) (*template.Template, error) {
	if tFname != "" {
		t, err := template.ParseFiles(tFname)
		if err != nil {
			return nil, err
		}
		return t, nil
	}
	return nil, nil
}

func parseContent(input []byte, tFname string) ([]byte, error) {
	output := blackfriday.Run(input)
	body := bluemonday.UGCPolicy().SanitizeBytes(output)

	t, err := template.New("mdp_").Parse(defaultTemplate)
	if err != nil {
		return nil, err
	}

	t, err = checkConfigurableTemplate(*t, tFname)
	if err != nil {
		return nil, err
	}

	c := content{
		Body:  template.HTML(body),
		Title: "Markdown Preview Tool",
	}
	//create a buffer to write to a file
	var buffer bytes.Buffer

	if err := t.Execute(&buffer, c); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func saveHTML(outFname string, data []byte) error {
	return os.WriteFile(outFname, data, 0644)
}

func preview(fname string) error {
	cName := ""
	cParams := []string{}

	switch runtime.GOOS {
	case "linux":
		cName = "xdg-open"
	case "windows":
		cName = "cmd.exe"
		cParams = []string{"/C", "start"}
	case "darwin":
		cName = "open"
	default:
		return fmt.Errorf("OS not supported")
	}
	cParams = append(cParams, fname)
	cPath, err := exec.LookPath(cName)
	if err != nil {
		return err
	}

	err = exec.Command(cPath, cParams...).Run()
	// quickfix to avoid the race condition opening the file before exclusion
	time.Sleep(2 * time.Second)
	return err
}
