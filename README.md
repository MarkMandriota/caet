# Cat Fetcher

## Installing
```bash
go install github.com/MarkMandriota/caet
```

## Usage
Get help:
```bash
caet -h
```

Common example:
```bash
caet
```
By default, will be installed 9 cats to `./cats`.

Flags:
```bash
  -N int
        number of workers (default 4)
  -d string
        cats destination directory (default "cats")
  -n int
        number of cats to fetch (default 9)
```

Example:
```bash
echo "https://api.thecatapi.com/v1/images/search" | caet -d cat_pics -n 96 -N 8
```

# TODO
* Add header key value pairs support