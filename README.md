# TinyPhotos

A small program that compresses JPG/JPEG images with metadata preservation with single image and bulk processing modes.

It uses [Tinify](https://tinypng.com/) to perform the processing and [ExifTool](https://exiftool.org/) to preserve the metadata.


## Why

Tinify's product TinyJPG automatically determines the best compression settings for your images with no noticeable differences at eyesight, in conjunction with ExifTool, this tool provides an easy way to compress your images and still preserve all the metadata possible.

## Requirements

- A Tinify API Key, that can be obtained on their [Developers page](https://tinypng.com/developers)
- ExifTool installed in your system and available on your $PATH as `exiftool`

## Setup

- Make a copy of the `.env.sample` file into a file named `.env` and fill out the `TINIFY_API_KEY` field with your Tinify API Key, you can get a developer key here: https://tinypng.com/developers

## Run the binary

```bash
tinyphotos <options>
```

## Run as developer

You can use the following command to quickly run the program without generating a binary:

```bash
go run . <options>
```

### Options / Flags
```
-file            <filepath>   Compresses a single file providing the relative or absolute filepath
-bulkfromfolder  <folderpath> Compresses all the files in a folder providing the relative or absolute path to the folder
-maxroutines     <number>     Maximum number of images that will be processed concurrently
```

## Build

### Build for your machine

Create an executable binary for your machine using the following command:

```bash
make build
```

### Build for all platforms

```bash
make build-all
```

### Roadmap
- [x] Concurrent execution for bulk processing mode
- [ ] Retry after failing
- [ ] Write log file (activate via CLI flag)

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License
MIT
