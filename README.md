# 1. Table of Contents
<!-- TOC -->

- [1. Table of Contents](#1-table-of-contents)
- [2. Installing](#2-installing)
- [3. Usage](#3-usage)
    - [3.1. Usage without arguments](#31-usage-without-arguments)

<!-- /TOC -->


# 2. Installing
``go get -U github.com/Necroforger/boorudl``

# 3. Usage
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

## 3.1. Usage without arguments
If executed without arguments, you will be asked to enter the information.
