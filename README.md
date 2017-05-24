# Boorudl

<!-- TOC -->

- [Boorudl](#boorudl)
    - [Installing](#installing)
    - [Usage](#usage)
        - [Usage without arguments](#usage-without-arguments)
        - [Excluding specific ratings](#excluding-specific-ratings)
    - [Other sites](#other-sites)

<!-- /TOC -->


## Installing

``go get -U github.com/Necroforger/boorudl``

Or download a version from the [releases](https://github.com/Necroforger/boorudl/releases).


## Usage
``boorudl https://safebooru.org -o "Cirno Pictures" -l 9 -t "Cirno blue"``

This will download a maximum of 9 images from safebooru with the tags *Cirno* and *Blue* to the directory *Cirno Pictures*.
The images will be named in the format [id].[extension].


| Flag              | Description                                             |
|-------------------|---------------------------------------------------------|
| -o                | Output directory for downloaded files                   |
| -p                | Page number                                             |
| -l                | Limit, or number of images to download                  |
| -t                | Space separated tags to search for                      |
| -r                | Get a random result (danbooru only)                     |
| --byscore         | Sort the results by score in order of highest to lowest |
| --exclude-ratings | Excludes posts with the supplied ratings.               |

### Usage without arguments
If executed without arguments, you will be asked to enter the information one line at a time.


### Excluding specific ratings
To exclude specific ratings from your search results you need to use the ``--exclude-ratings`` flag which accepts a non space separated list of ratings.

| Rating | type         |
|--------|--------------|
| q      | Questionable |
| e      | Explicit     |
| s      | Safe         |


Example:

``boorudl https://safebooru.org -o "Cirno Pictures" -l 9 -t "Cirno blue" --exclude-ratings qe``

Will exclude all ratings that are questionable and explicit.

## Other sites

| Site       | Info                                         |
|------------|----------------------------------------------|
| google.com | Download up to 100 images from google images |
