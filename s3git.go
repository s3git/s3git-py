/*
 * Copyright 2016 Frank Wessels <fwessels@xs4all.nl>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

// #cgo pkg-config: python3
// #define Py_LIMITED_API
// #include <Python.h>
import "C"

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/s3git/s3git-go"
)

//export s3git_init_repository
func s3git_init_repository(path *C.char) int {

	_, err := s3git.InitRepository(C.GoString(path))
	if err != nil {
		return -1
	}

	return 0
}

//export s3git_open_repository
func s3git_open_repository(path *C.char) int {

	_, err := s3git.OpenRepository(C.GoString(path))
	if err != nil {
		return -1
	}

	return 0
}

//export s3git_clone
func s3git_clone(url, path, accessKey, secretKey *C.char) int {

	options := []s3git.CloneOptions{}
	if access := C.GoString(accessKey); access != "" {
		options = append(options, s3git.CloneOptionSetAccessKey(access))
	}
	if secret := C.GoString(secretKey); secret != "" {
		options = append(options, s3git.CloneOptionSetSecretKey(secret))
	}

	_, err := s3git.Clone(C.GoString(url), C.GoString(path), options...)
	if err != nil {
		return -1
	}

	return 0
}

//export s3git_add
func s3git_add(path, filename *C.char) *C.char {

	repo, err := s3git.OpenRepository(C.GoString(path))
	if err != nil {
		return C.CString("")
	}

	file, err := os.Open(C.GoString(filename))
	if err != nil {
		return C.CString("")
	}

	key, _, err := repo.Add(file)
	if err != nil {
		return C.CString("")
	}

	return C.CString(key)
}

//export s3git_commit
func s3git_commit(path, message *C.char) int {

	repo, err := s3git.OpenRepository(C.GoString(path))
	if err != nil {
		return -1
	}

	repo.Commit(C.GoString(message))

	return 0
}

//export s3git_get
func s3git_get(path, hash *C.char) *C.char {

	repo, err := s3git.OpenRepository(C.GoString(path))
	if err != nil {
		return C.CString("")
	}

	r, err := repo.Get(C.GoString(hash))
	if err != nil {
		return C.CString("")
	}

	// TODO: Return stream directly instead of using temp file (and prevent dangling file)
	tmpfile, err := ioutil.TempFile("", "s3git")
	if err != nil {
		return C.CString("")
	}

	io.Copy(tmpfile, r)

	if err := tmpfile.Close(); err != nil {
		return C.CString("")
	}

	return C.CString(tmpfile.Name())
}

//export s3git_push
func s3git_push(path *C.char, hydrate bool) int {

	repo, err := s3git.OpenRepository(C.GoString(path))
	if err != nil {
		return -1
	}

	repo.Push(hydrate, func(total int64) {})

	return 0
}

//export s3git_pull
func s3git_pull(path *C.char) int {

	repo, err := s3git.OpenRepository(C.GoString(path))
	if err != nil {
		return -1
	}

	repo.Pull(func(total int64) {})

	return 0
}

//export s3git_list
func s3git_list(path, hash *C.char) *C.char {

	repo, err := s3git.OpenRepository(C.GoString(path))
	if err != nil {
		return C.CString("")
	}

	list, _ := repo.List(C.GoString(hash))

	response := []string{}

	count := 0
	for l := range list {
		response = append(response, l)

		count++
		if count > 1000 {
			break
		}
	}

	return C.CString(strings.Join(response, ","))
}

//export s3git_list_commits
func s3git_list_commits(path *C.char) *C.char {

	repo, err := s3git.OpenRepository(C.GoString(path))
	if err != nil {
		return C.CString("")
	}

	list, _ := repo.ListCommits("")

	result := []s3git.Commit{}

	count := 0
	for l := range list {
		result = append(result, l)

		count++
		if count > 1000 {
			break
		}
	}

	b, err := json.Marshal(result)

	return C.CString(string(b))
}

func main() {}
