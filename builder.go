package main

import (
	"os"
	"strings"

	"google.golang.org/protobuf/proto"

	router "github.com/v2fly/v2ray-core/v5/app/router/routercommon"
)

func main() {
	data, err := os.ReadFile("base/geosite.dat")
	if err != nil {
		panic(err)
	}

	var geo router.GeoSiteList
	if err := proto.Unmarshal(data, &geo); err != nil {
		panic(err)
	}

	// читаем кастом
	files, err := os.ReadDir("data")
	if err != nil {
		panic(err)
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		content, err := os.ReadFile("data/" + f.Name())
		if err != nil {
			continue
		}

		lines := strings.Split(string(content), "\n")

		var domains []*router.Domain

		for _, l := range lines {
			l = strings.TrimSpace(l)
			if l == "" {
				continue
			}

			d := &router.Domain{}

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
				d.Value = l
			}

			domains = append(domains, d)
		}

		name := strings.ToLower(strings.TrimSuffix(f.Name(), ".txt"))

		target := ""

		// =========================
		// CATEGORY MAP (твои правила)
		// =========================
		switch {
		case strings.Contains(name, "direct"):
			target = "private"

		case strings.Contains(name, "ads"):
			target = "ads"

		case strings.Contains(name, "apple"):
			target = "apple"

		case strings.Contains(name, "tm"):
			target = "tm-rules"

		default:
			target = "proxy"
		}

		// =========================
		// FIX: visibility in PassWall2 Geo View
		// =========================
		geo.Entry = append(geo.Entry, &router.GeoSite{
			CountryCode: "geosite:" + target,
			Domain:      domains,
		})
	}

	out, err := proto.Marshal(&geo)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("geosite.dat", out, 0644)
	if err != nil {
		panic(err)
	}
}
