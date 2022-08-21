# Cat Fetcher

## Installing
```bash
go install github.com/MarkMandriota/caet/cmd/caet
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
  -p string
        random cats providers splited by ";" (default "https://thiscatdoesnotexist.com/;https://api.thecatapi.com/v1/images/search")
```

Example:
```bash
caet -d cat_pics -p "https://api.thecatapi.com/v1/images/search" -n 96 -N 8
```