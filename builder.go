package main

import (
	"os"
	"strings"

	"google.golang.org/protobuf/proto"

	router "github.com/v2fly/v2ray-core/v5/app/router/routercommon"
)

func main() {
	data, _ := os.ReadFile("base/geosite.dat")

	var geo router.GeoSiteList
	proto.Unmarshal(data, &geo)

	// читаем кастом
	files, _ := os.ReadDir("data")

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		content, _ := os.ReadFile("data/" + f.Name())
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
			} else {
				d.Type = router.Domain_RootDomain
				d.Value = l
			}

			domains = append(domains, d)
		}

		name := strings.ToLower(strings.TrimSuffix(f.Name(), ".txt"))

		target := ""

		// 🔥 МАППИНГ В СТАНДАРТ (КЛЮЧЕВОЕ)
		switch {
		case strings.Contains(name, "direct"):
			target = "private"
		case strings.Contains(name, "ads"):
			target = "ads"
		case strings.Contains(name, "apple"):
			target = "apple"
		default:
			target = "proxy"
		}

		geo.Entry = append(geo.Entry, &router.GeoSite{
			CountryCode: target,
			Domain:      domains,
		})
	}

	out, _ := proto.Marshal(&geo)
	os.WriteFile("geosite.dat", out, 0644)
}
