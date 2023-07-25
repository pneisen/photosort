# photosort
Quick and dirty golang to sort my photos exported from OSX Photos. Traverses a directory tree and
uses exif to copy photos to a flat directory structure grouped by the year and month the photo was taken (ie 2023-7).
It won't copy over same named files and prints those for manual move. If it encounters a file that doesn't
have exif, it copies the file to an unsorted directory. The intention of this sorting is for eventual upload into S3 buckets. 

Usage:
  - Clone repo
  - go mod tidy
  - go build
  - `./photosort <from dir> <to dir>`