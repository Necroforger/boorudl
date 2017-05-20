package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Necroforger/Boorudl/downloader"
)

var flagset = flag.NewFlagSet("Boorudl", flag.ExitOnError)

// Command line arguments
var (
	OutputDir string
	BooruURL  string
	Tags      string
	Limit     int
	Page      int
	Random    bool
)

// LogError logs errors if they are not nil
func LogError(e error) {
	if e != nil {
		fmt.Println(e)
	}
}

// ParseFlags parses the commandline flags
func ParseFlags() {
	flagset.StringVar(&OutputDir, "o", "", "Output directory for downloaded files")
	flagset.IntVar(&Page, "p", 1, "Page to start downloading from.")
	flagset.IntVar(&Limit, "l", 1, "Maximum number of images to download")
	flagset.StringVar(&Tags, "t", "", "Space separated tags to search for")
	flagset.BoolVar(&Random, "r", false, "Specifies if the result should be random. Only works on danbooru")

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "usage: boorudl [booru site] [flags]")
		flagset.PrintDefaults()
		os.Exit(1)
	}

	BooruURL = os.Args[1]

	flagset.Parse(os.Args[2:])
	fmt.Println("Limit: ", Limit)
}

func main() {
	ParseFlags()

	if OutputDir == "" {
		fmt.Println("Invalid output directory")
		return
	}

	results, err := downloader.Search(BooruURL, downloader.SearchQuery{
		Tags:   Tags,
		Limit:  Limit,
		Page:   Page,
		Random: Random,
	})
	if err != nil {
		fmt.Println("Error obtaining information from booru: ", err)
	}

	sort.Sort(downloader.SearchResultsByScore(results))

	os.MkdirAll(OutputDir, 0666)
	for i, v := range results {
		LogError(SaveFileFromURL(v.ImageURL, filepath.Join(OutputDir, fmt.Sprint(v.ID))))
		fmt.Printf("%d/%d\n", i+1, len(results))
	}
}

// SaveFileFromURL will download the file from the specified url,
// Infer the content type, and save with the correct extension.
// 		URL:		The url to retrieve the file from
//		filename:	the path to save the file to excluding the file extenstion
// 					as it will be inferred.
func SaveFileFromURL(URL string, path string) error {

	resp, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Infer the file extension from the response body.
	sampleBytes, extension, err := DetectContentTypeFromReader(resp.Body)
	if err != nil {
		return err
	}

	// Check that the extension is of a known file type.
	if extension == "octet-stream" {
		return errors.New("invalid content type")
	}

	dlfile, err := os.OpenFile(path+"."+extension, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	defer dlfile.Close()

	// Prepend the sample data that was taken from the reader to infer the content type.
	_, err = dlfile.Write(sampleBytes)
	if err != nil {
		return err
	}

	// Copy the rest of the response into the file
	_, err = io.Copy(dlfile, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// DetectContentTypeFromReader detects the content type from an io.Reader stream
// It will return the bytes used to infer the content type.
func DetectContentTypeFromReader(reader io.Reader) ([]byte, string, error) {
	sampleData := make([]byte, 512)
	numread, err := reader.Read(sampleData)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return nil, "", err
	}

	extension := strings.Split(http.DetectContentType(sampleData), "/")[1]

	return sampleData[:numread], extension, nil
}
