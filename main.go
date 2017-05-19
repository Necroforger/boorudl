package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/Necroforger/Boorudl/downloader"
)

func main() {
	boo := downloader.NewDanbooru()
	results, err := boo.Search(downloader.SearchQuery{
		Tags:  "clownpiece",
		Limit: 2000,
		Page:  1,
	})
	if err != nil {
		fmt.Println(err)
	}

	os.MkdirAll("BooruSearch", 0666)
	for i, v := range results {
		saveImage(v.ImageURL, "BooruSearch/"+fmt.Sprint(v.ID))
		fmt.Printf("%d/%d\n", i, len(results)-1)
	}
}

func saveImage(url string, path string) {

	res, err := http.Get(url)
	if err != nil {
		fmt.Println("error fetching thumbnail: ", err)
	}

	defer res.Body.Close()
	dbuf := make([]byte, 256)
	numread, err := res.Body.Read(dbuf)
	if err != nil && (err != io.EOF || err != io.ErrUnexpectedEOF) {
		fmt.Println("error reading: ", err)
		return
	}

	ext := strings.Split(http.DetectContentType(dbuf), "/")[1]

	// No image content detected
	if ext == "octet-stream" {
		return
	}

	fmt.Println("content-type: ", ext)

	file, err := os.OpenFile(path+"."+ext, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("error creating file: ", err)
		return
	}
	defer file.Close()

	_, err = file.Write(dbuf[:numread])
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = io.Copy(file, res.Body)
	if err != nil {
		fmt.Println("error writing image: ", err)
		return
	}
}
