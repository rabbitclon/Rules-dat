package main

import (
	"fmt"
	"os"
	"strings"

	"google.golang.org/protobuf/proto"

	router "github.com/v2fly/v2ray-core/v5/app/router/routercommon"
)

func main() {
	baseFile := "build/base.dat"

	data, err := os.ReadFile(baseFile)
	if err != nil {
		panic(err)
	}

	var base router.GeoSiteList
	if err := proto.Unmarshal(data, &base); err != nil {
		panic(err)
	}

	// читаем твои кастомные правила
	customDir := "build/custom"

	files, _ := os.ReadDir(customDir)

	for _, f := range files {
		name := strings.ToLower(strings.TrimSuffix(f.Name(), ".txt"))

		content, _ := os.ReadFile(customDir + "/" + f.Name())
		lines := strings.Split(string(content), "\n")

		var domains []*router.Domain

		for _, l := range lines {
			l = strings.TrimSpace(l)
			if l == "" {
				continue
			}

			d := &router.Domain{
				Value: l,
			}

			// если full:
			if strings.HasPrefix(l, "full:") {
				d.Type = router.Domain_Full
				d.Value = strings.TrimPrefix(l, "full:")
			} else if strings.HasPrefix(l, "keyword:") {
				d.Type = router.Domain_Plain
				d.Value = strings.TrimPrefix(l, "keyword:")
			} else if strings.HasPrefix(l, "regexp:") {
				d.Type = router.Domain_Regex
				d.Value = strings.TrimPrefix(l, "regexp:")
			} else {
				d.Type = router.Domain_RootDomain
			}

			domains = append(domains, d)
		}

		// добавляем новую категорию в base
		base.Entry = append(base.Entry, &router.GeoSite{
			CountryCode: name,
			Domain:      domains,
		})

		fmt.Println("Added category:", name)
	}

	out, err := proto.Marshal(&base)
	if err != nil {
		panic(err)
	}

	os.WriteFile("geosite.dat", out, 0644)

	fmt.Println("Build complete: geosite.dat")
}
