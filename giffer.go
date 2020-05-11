/*

Giffer.go

Copyright (C) 2021 the Free Software Foundation. All Rights Reserved.
Based on Giffy (copyright 2013 Google Inc.)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
     http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Giffer reads all the JPEG and PNG files from the current directory
(the directory that it is in) and writes them to an animated GIF as "out.gif".

The animation is 130ms; you may edit that below.
*/

package main

import (
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "image/jpeg"
	_ "image/png"
)

func main() {

	log.Printf("Giffy v1.1\n")
	log.Printf("All Rights Reserved.\n\n")
	log.Printf("Giffy takes all the images in the folder that\n")
	log.Printf("it is in, and turns them into an animated GIF.\n\n")

	fs, err := dirFiles(".")
	if err != nil {
		log.Fatal("Error reading files in this directory:", err)
	}
	var ms []*image.Paletted
	for i, n := range fs {
		log.Printf("Reading %v [%d/%d]\n", n, i+1, len(fs))
		m, err := readImage(n)
		if err != nil {
			log.Fatalf("Error reading image: %v: %v", n, err)
		}
		r := m.Bounds()
		pm := image.NewPaletted(r, palette.Plan9)
		draw.FloydSteinberg.Draw(pm, r, m, image.ZP)
		ms = append(ms, pm)
	}
	ds := make([]int, len(ms))
	for i := range ds {
		ds[i] = 130
	}
	const out = "out.gif"
	log.Println("Generating ", out)
	f, err := os.Create(out)
	if err != nil {	log.Fatalf("Error creating %v: %v", out, err) }
	err = gif.EncodeAll(f, &gif.GIF{Image: ms, Delay: ds, LoopCount: 999999999})
	if err != nil { log.Fatalf("Error writing %v: %v", out, err) }
	err = f.Close()
	if err != nil { log.Fatalf("Error closing %v: %v", out, err) }
}

var validExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
}

func dirFiles(dir string) (names []string, err error) {
	fs, err := ioutil.ReadDir(".")
	if err != nil {
		return nil, err
	}
	for _, fi := range fs {
		n := fi.Name()
		if !validExtensions[filepath.Ext(n)] {
			continue
		}
		names = append(names, n)
	}
	sort.Sort(filenames(names))
	return
}

type filenames []string

func (s filenames) Len() int      { return len(s) }
func (s filenames) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s filenames) Less(i, j int) bool {
	if filepath.Ext(s[i]) == filepath.Ext(s[j]) {
		a, b := stripExt(s[i]), stripExt(s[j])
		if (strings.HasPrefix(a, b) || strings.HasPrefix(b, a)) &&
			strings.Contains(a, "#") != strings.Contains(b, "#") {
			return strings.Contains(b, "#")
		}
	}
	return s[i] < s[j]
}

func stripExt(name string) string {
	return strings.TrimSuffix(name, filepath.Ext(name))
}

func readImage(name string) (image.Image, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	m, _, err := image.Decode(f)
	return m, err
}
