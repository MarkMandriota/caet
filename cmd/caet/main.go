// Copyright 2022 Mark Mandriota
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/MarkMandriota/caet"
)

var (
	dstD = flag.String(`d`, `cats`, `cats destination directory`)
	catN = flag.Int(`n`, 9, `number of cats to fetch`)
	proN = flag.Int(`p`, runtime.NumCPU(), `number of processes`)
)

func init() {
	flag.Parse()
}

func main() {
	cats := make(chan *caet.CatFile, *catN)

	fetcher := caet.NewFetcher()

	if err := fetcher.SR.Load(os.Stdin); err != nil {
		log.Fatalln(err)
	}

	ctx := context.Background()

	for i := *proN; i > 0; i-- {
		go fetcher.Run(ctx, cats)
	}

	if err := os.Mkdir(*dstD, 0777); err != nil && !errors.Is(err, os.ErrExist) {
		log.Fatal(err)
	}

	for i := 0; i < *catN; i++ {
		catF := <-cats

		path := fmt.Sprintf("%s/%d.%s", *dstD, i+1, catF.Kind)

		file, err := os.Create(path)
		if err != nil {
			log.Printf("error creating file %s: %v\n", path, err)
		}

		if _, err := file.Write(catF.Body); err != nil {
			log.Printf("error writing file %s: %v\n", path, err)
		}

		file.Close()

		fmt.Printf("\rfetched: %d", i+1)
	}
}
