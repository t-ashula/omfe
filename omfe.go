package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-ini/ini"
	"gopkg.in/yaml.v2"
)

func main() {
	if len(os.Args) != 2 {
		usage(filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	fname := os.Args[1]
	contents, err := loadIndexIni(fname)
	if err != nil {
		die("load ini faild", err)
	}

	feeds, err := extractFeedList(contents)
	if err != nil {
		die("extract feed failed", err)
	}

	err = writeFeedList(feeds)
	if err != nil {
		die("write feed failed", err)
	}
}

type feedInfo struct {
	URL   string
	Title string
}

func writeFeedList(feeds []feedInfo) error {
	data, err := yaml.Marshal(feeds)
	if err != nil {
		return err
	}
	fmt.Print(string(data))
	return nil
}

func extractFeedList(contents []byte) ([]feedInfo, error) {
	cfg, err := ini.LoadSources(ini.LoadOptions{AllowBooleanKeys: true}, contents)
	if err != nil {
		die(fmt.Sprintf("parse ini failed %s", contents[:20]), err)
	}

	sections := cfg.Sections()
	var feeds []feedInfo
	for _, s := range sections {
		if !strings.HasPrefix(s.Name(), "Index ") {
			continue
		}
		k, err := s.GetKey("Search Text")
		if k == nil || err != nil {
			continue
		}
		v := k.Value()
		if !strings.HasPrefix(v, "http") {
			continue
		}
		t := s.Key("Name").Value()
		feeds = append(feeds, feedInfo{URL: v, Title: t})
	}
	return feeds, nil
}

func loadIndexIni(fname string) ([]byte, error) {
	content, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, err
	}

	index := bytes.Index(content, []byte("Opera Preferences"))
	if !(index == 0 || index == 3) {
		return nil, fmt.Errorf("not opera preference file")
	}
	index = bytes.Index(content, []byte("[Indexer]"))
	if index == -1 {
		return nil, fmt.Errorf("not index.init file")
	}

	return content[index:], nil
}

func die(msg string, err error) {
	fmt.Fprintf(os.Stderr, "%s, %s\n", msg, err)
	os.Exit(1)
}

func usage(name string) {
	fmt.Printf("usage: %s /path/to/operamail/index.ini\n", name)
}
