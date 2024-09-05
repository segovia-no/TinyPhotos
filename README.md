# TinyPhotos

A small program that compresses JPG/JPEG images with metadata preservation with single image and bulk processing modes.

It uses [Tinify](https://tinypng.com/) to perform the processing and [ExifTool](https://exiftool.org/) to preserve the metadata.


## Why

Tinify's product TinyJPG automatically determines the best compression settings for your images with no noticeable differences at eyesight, in conjunction with ExifTool, this tool provides an easy way to compress your images and still preserve all the metadata possible.

Compared to the web version, using the API directly gives you access to upload files more than 5MB

## Requirements

- A Tinify API Key, that can be obtained on their [Developers page](https://tinypng.com/developers)
- ExifTool installed in your system and available on your $PATH as `exiftool`

## Setup

You need to setup your Tinify API key developer, this program allows you to set it up using either an `.env` file or the `apikey` flag:

A) __Use of .env file__: Make a copy of the `.env.sample` file into a file named `.env` and fill out the `TINIFY_API_KEY` field with your Tinify API Key.

B) __Use of -apikey flag__: When running the binary, you can pass the api key directly using the `-apikey` flag.

>Note: If both options are provided, the `-apikey` flag will take priority over the `.env` file.

## Run the binary

```bash
./tinyphotos <options>
```

## Run as developer

You can use the following command to quickly run the program without generating a binary:

```bash
go run . <options>
```

## Options / Flags
```
-apikey          <string>               An alternative method to the .env file, it allows the user to provide the api key directly as a flag
-file            <filepath>             Compresses a single file providing the relative or absolute filepath
-bulkfromfolder  <folderpath>           Compresses all the files in a folder providing the relative or absolute path to the folder
-maxroutines     <number (default: 4)>  Maximum number of images that will be processed concurrently
-log                                    If provided, the output of the program will be saved into a file
-maxretries      <number (default: 2)>  Max number of retries if a request to the Tinify API fails per file
```

>Note: Adjust the `maxroutines` option according to your network's upload speed, a higher value might generate more failure rates while making compression requests to the Tinify API.
>Note: The `maxretries` option only works when using the `bulkfromfolder` option.

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
- [x] Support API key passthrough using `-apikey` flag
- [x] Concurrent execution for bulk processing mode
- [x] Retry after failing with a `-maxretries` flag
- [x] Write log file (activate via CLI flag)
- [ ] Limit amount of compressions made to the Tinify API via `-complimit` flag. Useful for controlling costs
- [ ] Investigate how to make the requests to the Tinify API faster for larger files

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License
MIT
