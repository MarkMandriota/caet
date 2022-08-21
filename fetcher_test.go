package caet

import (
	"bytes"
	_ "embed"
	"log"
	"sync"

	"testing"

	http "github.com/valyala/fasthttp"
)

//go:embed example/1.jpeg
var catRaw []byte

const (
	addr = "localhost:9090"

	dPath = "/"

	iPath = "/1.jpeg"

	mess = `Please, follow this URL to get image: 
` + `http://` + addr + iPath + `
Thanks.`
)

func init() {
	go func() {
		if err := http.ListenAndServe(addr, func(ctx *http.RequestCtx) {
			switch string(ctx.Path()) {
			case dPath:
				ctx.Response.Header.Set("Content-Type", "text/plain")
				if _, err := ctx.WriteString(mess); err != nil {
					log.Println(err)
				}
			case iPath:
				ctx.Response.Header.Set("Content-Type", "image/jpeg")
				if _, err := ctx.Write(catRaw); err != nil {
					log.Println(err)
				}
			}
		}); err != nil {
			log.Fatalln(err)
		}
	}()
}

func TestFetcherFetchNewer(t *testing.T) {
	tests := []struct {
		name string
		addr string
	}{
		{
			name: "direct API",
			addr: "http://" + addr + iPath,
		},
		{
			name: "indirect API",
			addr: "http://" + addr + dPath,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fetcher := NewFetcher()
			fetcher.SR.SetList([]string{test.addr})

			if body, kind := fetcher.FetchNewer(); !(bytes.Equal(kind, []byte("jpeg")) && bytes.Equal(catRaw, body)) {
				t.Fatal("images do not match")
			}
		})
	}
}

func BenchmarkFetcherFetchNewer(b *testing.B) {
	fetcher := NewFetcher()
	fetcher.SR.SetList([]string{"http://" + addr})

	for i := 0; i < b.N; i++ {
		fetcher.FetchNewer()
		fetcher.hashes = sync.Map{}
	}
}
