# boorudl
Installing:

``go get -U github.com/Necroforger/boorudl``

# Usage
``boorudl https://safebooru.org -o "Cirno Pictures" -l 9 -t "Cirno blue"``

This will download a maximum of 9 images from safebooru with the tags *Cirno* and *Blue* to the directory *Cirno Pictures*.
The images will be named in the format [id].[extension].

| Flag | Description |
|------:|-----------:|
|-o      |Output directory for downloaded files|
|-p      | Page number|
|-l      | Limit, or number of images to download|
|-t      | Space separated tags to search for|
|-r      | Get a random result (danbooru only)|
