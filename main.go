package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Necroforger/Boorudl/downloader"
	httpclient "github.com/mreiferson/go-httpclient"
)

var flagset = flag.NewFlagSet("Boorudl", flag.ExitOnError)

// Create an http client with a five second timeout on reads.
// The timer will reset every time data is received.
var client = http.Client{Transport: &httpclient.Transport{
	ConnectTimeout:        time.Second * 10,
	ReadWriteTimeout:      time.Second * 5,
	ResponseHeaderTimeout: time.Second * 10,
}}

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
	flagset.IntVar(&Page, "p", 0, "Page to start downloading from.")
	flagset.IntVar(&Limit, "l", 1, "Maximum number of images to download")
	flagset.StringVar(&Tags, "t", "", "Space separated tags to search for")
	flagset.BoolVar(&Random, "r", false, "Specifies if the result should be random. Only works on danbooru")

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "usage: boorudl [booru site] [flags]")
		flagset.PrintDefaults()
		SetFromUserInput()
		return
	}

	BooruURL = os.Args[1]

	flagset.Parse(os.Args[2:])
}

// RequestInput requests input from the user
func RequestInput(query string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(query)
	fmt.Print(">")

	res, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	res = strings.Replace(res, "\n", "", -1)
	res = strings.Replace(res, "\r", "", -1)

	return res, nil
}

// SetFromUserInput asks the user to enter fields if no arguments to boorudl have been provided
func SetFromUserInput() {
	// var err error

	BooruURL, _ = RequestInput("Booru URL: ")
	Tags, _ = RequestInput("Tags: ")

	limitstr, err := RequestInput("limit: default(1)")
	if err == nil {
		if n, err := strconv.Atoi(limitstr); err == nil {
			Limit = n
		} else {
			Limit = 1
		}
	}

	pagestr, err := RequestInput("page number: default(0)")
	if err == nil {
		if n, err := strconv.Atoi(pagestr); err == nil {
			Page = n
		} else {
			Page = 0
		}
	}

	OutputDir, _ = RequestInput("Output directory: (default is the current console directory")
}

func main() {
	ParseFlags()

	results, err := downloader.Search(BooruURL, downloader.SearchQuery{
		Tags:   Tags,
		Limit:  Limit,
		Page:   Page,
		Random: Random,
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	// If not enough images have been found, Search the next page until 'Limit' results have been found
	// Or nothing is returned.
	for pagenum := Page + 1; len(results) < Limit; pagenum++ {

		r, err := downloader.Search(BooruURL, downloader.SearchQuery{
			Tags:   Tags,
			Limit:  Limit,
			Page:   pagenum,
			Random: Random,
		})
		if err != nil {
			fmt.Println("Finished searching for images: ", err)
			break
		}
		results = append(results, r...)
		fmt.Println("Added ", len(r), "images to queue")

	}
	fmt.Println("found ", len(results), "images")

	if OutputDir != "" {
		err = os.MkdirAll(OutputDir, 0600)
		if err != nil {
			fmt.Println("Error creating output directory", err)
			return
		}
	}

	fmt.Println("Attempting to save images...")

	for i, v := range results {
		LogError(SaveFileFromURL(v.ImageURL, filepath.Join(OutputDir, fmt.Sprint(v.ID))))
		fmt.Printf("%d/%d\t%s\n", i+1, len(results), v.ImageURL)
	}

}

// SaveFileFromURL will download the file from the specified url,
// Infer the content type, and save with the correct extension.
// 		URL:		The url to retrieve the file from
//		filename:	the path to save the file to excluding the file extenstion
// 					as it will be inferred.
func SaveFileFromURL(URL string, path string) error {

	resp, err := client.Get(URL)
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
