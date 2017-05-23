package main

import (
	"fmt"
	"os"

	"github.com/Necroforger/boorudl/downloader"
	astilectron "github.com/Necroforger/go-astilectron"
)

var (
	workingDirectory string
	electron         *astilectron.Astilectron
	window           *astilectron.Window
)

var posts downloader.Posts

func main() {

	var err error
	workingDirectory, err = os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}

	// TODO Create application icons.
	if electron, err = astilectron.New(astilectron.Options{
		AppName:            "Boorudl",
		AppIconDarwinPath:  "",
		AppIconDefaultPath: "",
		BaseDirectoryPath:  workingDirectory,
	}); err != nil {
		fmt.Println(err)
		return
	}
	defer electron.Close()
	electron.HandleSignals()

	// Attempt to start the electron session.
	// If electron does not exist, download it.
	if err = electron.Start(); err != nil {
		fmt.Println(err)
		return
	}

	if window, err = electron.NewWindow(workingDirectory+"/index.html", &astilectron.WindowOptions{
		Width:     astilectron.PtrInt(800),
		Height:    astilectron.PtrInt(600),
		Center:    astilectron.PtrBool(true),
		Resizable: astilectron.PtrBool(false),
	}); err != nil {
		fmt.Println(err)
		return
	}
	if err = window.Create(); err != nil {
		fmt.Println(err)
		return
	}

	electron.Wait()
}
