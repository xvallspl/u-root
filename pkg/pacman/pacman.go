// Package pacman provides a hacked together API for Arch Linux packages.
package pacman

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type pkgInfo struct {
	Name      string   `json:"pkgname"`
	Arch      string   `json:"arch"`
	Repo      string   `json:"repo"`
	Desc      string   `json:"pkgdesc"`
	Deps      []string `json:"depends"`
	Licenses  []string `json:"licenses"`
	Filename  string   `json:"filename"`
	Conflicts []string `json:"conflicts"`
	Version   string   `json:"pkgver"`
}

type Pkg struct {
	info *pkgInfo
}

const baseURL = "https://www.archlinux.org/packages/"

func pkgURL(repo, arch, name string) string {
	return fmt.Sprintf("%s/%s/%s/%s", baseURL, repo, arch, name)
}

func info(repo, arch, name string) (*pkgInfo, error) {
	resp, err := http.Get(fmt.Sprintf("%s/json", pkgURL(repo, arch, name)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var info pkgInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}
	return &info, nil
}

func Package(repo, arch, name string) (*Pkg, error) {
	i, err := info(repo, arch, name)
	if err != nil {
		return nil, err
	}
	return &Pkg{
		info: i,
	}, nil
}

func (p *Pkg) Install(chroot string) error {

}
