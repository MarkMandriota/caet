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
	"bytes"
	"context"
	"crypto/md5"
	"sync"

	http "github.com/valyala/fasthttp"
)

type CatFile struct {
	Body []byte
	Kind []byte
}

func NewCatFile(body []byte, kind []byte) *CatFile {
	return &CatFile{body, kind}
}

type Fetcher struct {
	SR SingularReferer

	client *http.Client

	hashes sync.Map
}

func NewFetcher() *Fetcher {
	return &Fetcher{
		client: &http.Client{},
	}
}

func (f *Fetcher) Run(ctx context.Context, cats chan<- *CatFile) {
	for {
		select {
		case cats <- NewCatFile(f.FetchNewer()):
		case <-ctx.Done():
			return
		}
	}
}

func (f *Fetcher) FetchNewer() (body []byte, kind []byte) {
	res := http.AcquireResponse()
	req := http.AcquireRequest()

next:
	req.SetRequestURI(f.SR.Next())
	for {
		if err := f.client.Do(req, res); err != nil {
			goto next
		}

		mode, kind := readContentType(res.Header.Peek("Content-Type"))

		body = res.Body()

		if bytes.Equal(mode, []byte("image")) {
			if _, ok := f.hashes.LoadOrStore(md5.Sum(body), struct{}{}); ok {
				goto next
			}

			http.ReleaseResponse(res)
			http.ReleaseRequest(req)
			return body, kind
		}

		req.SetRequestURIBytes(rWebReference.Find(body))
	}
}
