package ripsrc

import (
	"regexp"
	"sort"

	"gopkg.in/src-d/go-license-detector.v2/licensedb"
	"gopkg.in/src-d/go-license-detector.v2/licensedb/filer"
)

type memoryfiler struct {
	filename string
	buf      []byte
}

// ReadFile returns the contents of a file given it's path.
func (f *memoryfiler) ReadFile(path string) (content []byte, err error) {
	return f.buf, nil
}

// ReadDir lists a directory.
func (f *memoryfiler) ReadDir(path string) ([]filer.File, error) {
	return []filer.File{
		filer.File{
			Name:  f.filename,
			IsDir: false,
		},
	}, nil
}

// Close frees all the resources allocated by this Filer.
func (f *memoryfiler) Close() {
}

// License holds details about detected license
type License struct {
	Name       string  `json:"license"`
	Confidence float32 `json:"confidence"`
}

var licenses = regexp.MustCompile("(LICENSE|README)(\\.(md|txt))?$")

func possibleLicense(filename string) bool {
	return licenses.MatchString(filename)
}

const minConfidenceLevel float32 = 0.85

func detect(filename string, buf []byte) (*License, error) {
	mf := &memoryfiler{filename, buf}
	kv, err := licensedb.Detect(mf)
	if err != nil {
		if err == licensedb.ErrNoLicenseFound {
			return nil, nil
		}
		return nil, err
	}
	if len(kv) > 0 {
		matches := make([]License, 0)
		for k, v := range kv {
			matches = append(matches, License{k, v})
		}
		if len(matches) > 0 {
			sort.Slice(matches, func(i, j int) bool {
				return matches[i].Confidence > matches[j].Confidence
			})
		}
		if matches[0].Confidence >= minConfidenceLevel {
			return &matches[0], nil
		}
	}
	return nil, nil
}
