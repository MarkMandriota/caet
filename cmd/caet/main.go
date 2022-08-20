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
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/MarkMandriota/caet"
)

var (
	dstD = flag.String(`d`, `cats`, `cats destination directory`)
	catP = flag.String(`p`, `https://thiscatdoesnotexist.com/`, `random cats providers splited by ";"`)
	catN = flag.Int(`n`, 9, `number of cats to fetch`)
	numW = flag.Int(`N`, runtime.NumCPU(), `number of workers`)
)

func init() {
	flag.Parse()
}

func main() {
	cats := make(chan *caet.CatFile, *catN)

	fetcher := &caet.Fetcher{}

	if *catP != "" {
		fetcher.SR.SetList(strings.Split(*catP, ";"))
	} else {
		fetcher.SR.Load(os.Stdin)
	}

	for i := *numW; i > 0; i-- {
		go fetcher.Run(cats)
	}

	os.Mkdir(*dstD, 0777)

	for i := 0; i < *catN; i++ {
		catF := <-cats

		path := fmt.Sprintf("%s/%d.%s", *dstD, i+1, catF.Kind)

		file, err := os.Create(path)
		if err != nil {
			log.Printf("error creating file %s: %v\n", path, err)
		}

		file.Write(catF.Body)

		file.Close()
	}
}
