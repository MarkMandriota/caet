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

package caet

import (
	"crypto/md5"
	"io"
	"net/http"
	"sync"
)

type CatFile struct {
	Body []byte
	Kind string
}

func NewCatFile(body []byte, kind string) *CatFile {
	return &CatFile{body, kind}
}

type Fetcher struct {
	SR SingularReferer

	Client http.Client

	hashes sync.Map
}

func (f *Fetcher) Run(cats chan<- *CatFile) {
	for {
		cats <- NewCatFile(f.FetchNewer())
	}
}

func (f *Fetcher) FetchNewer() (body []byte, kind string) {
next:
	ref := f.SR.Next()
	for {
		resp, err := f.Client.Get(ref)
		if err != nil {
			goto next
		}

		text, err := io.ReadAll(resp.Body)
		if err != nil {
			continue
		}

		kind, typ := readContentType(resp.Header.Get("Content-Type"))

		if kind == "image" {
			if _, ok := f.hashes.LoadOrStore(md5.Sum(text), struct{}{}); ok {
				goto next
			}
			return text, typ
		}

		ref = string(rWebReference.Find(text))
	}
}
