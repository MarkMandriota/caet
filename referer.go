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
	"bufio"
	"io"
	"math/rand"
	"regexp"
	"sync"
	"time"
)

var rWebReference = regexp.MustCompile(`https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`)

type SingularReferer struct {
	m sync.Mutex

	list []string
}

func (r *SingularReferer) Load(cfg io.Reader) error {
	r.m.Lock()
	scanner := bufio.NewScanner(cfg)

	if cfg != nil {
		r.list = r.list[:0]
	}

	for scanner.Scan() {
		r.list = append(r.list, scanner.Text())
	}

	rand.Seed(time.Now().UnixNano())

	r.m.Unlock()
	return scanner.Err()
}

func (r *SingularReferer) Next() (ref string) {
	r.m.Lock()
	ref = r.list[rand.Intn(len(r.list))]

	r.m.Unlock()
	return ref
}
