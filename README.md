# 1. Boorudl

<!-- TOC -->

- [1. Boorudl](#1-boorudl)
    - [1.1. Installing](#11-installing)
    - [1.2. Usage](#12-usage)
        - [1.2.1. Usage without arguments](#121-usage-without-arguments)

<!-- /TOC -->


## 1.1. Installing

``go get -U github.com/Necroforger/boorudl``

Or download a version from the [releases](https://github.com/Necroforger/boorudl/releases).


## 1.2. Usage
``boorudl https://safebooru.org -o "Cirno Pictures" -l 9 -t "Cirno blue"``

This will download a maximum of 9 images from safebooru with the tags *Cirno* and *Blue* to the directory *Cirno Pictures*.
The images will be named in the format [id].[extension].


| Flag | Description                            |
|------|----------------------------------------|
| -o   | Output directory for downloaded files  |
| -p   | Page number                            |
| -l   | Limit, or number of images to download |
| -t   | Space separated tags to search for     |
| -r   | Get a random result (danbooru only)    |

### 1.2.1. Usage without arguments
If executed without arguments, you will be asked to enter the information one line at a time.
