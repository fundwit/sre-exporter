package localize

import (
	"sort"
	"strconv"
	"strings"
)

type AcceptLanguages []AcceptLanguage

type AcceptLanguage struct {
	Lang   string
	Weight float64
}

func (a AcceptLanguages) Len() int {
	return len(a)
}
func (a AcceptLanguages) Less(i, j int) bool {
	if a[i].Weight != a[j].Weight {
		return a[i].Weight > a[j].Weight
	}

	if len(a[i].Lang) != len(a[j].Lang) {
		return len(a[i].Lang) > len(a[j].Lang)
	}

	return a[i].Lang > a[j].Lang
}
func (a AcceptLanguages) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// ParseAcceptLanguage example: zh-CN,zh;q=0.8,en;q=0.6,zh-TW;q=0.4
func ParseAcceptLanguage(lang string) AcceptLanguages {
	parts := strings.Split(lang, ",")
	acceptLanguages := AcceptLanguages{}
	langs := map[string]float64{}
	var langsExt []string

	for _, p := range parts {
		var langName string
		var weight float64

		if strings.Contains(p, ";") {
			langWeightPair := strings.Split(p, ";")
			lang := langWeightPair[0]
			weightExpr := langWeightPair[1]
			langName = strings.TrimSpace(lang)

			if strings.Contains(weightExpr, "=") {
				qvPair := strings.Split(weightExpr, "=")
				q := qvPair[0]
				value := qvPair[1]
				if q == "q" {
					w, err := strconv.ParseFloat(value, 64)
					if err == nil {
						weight = w
					}
				}
			}
		} else {
			langName = p
			weight = 1
		}

		if strings.Contains(langName, "-") {
			langNameBase := strings.Split(langName, "-")[0]
			if langNameBase != "" {
				langsExt = append(langsExt, langNameBase)
			}
		}

		if langName == "*" {
			langName = ""
		}

		if langName != "" {
			langs[langName] = weight
		}
	}

	for _, l := range langsExt {
		if _, exist := langs[l]; !exist {
			langs[l] = 0
		}
	}

	for l, w := range langs {
		acceptLanguages = append(acceptLanguages, AcceptLanguage{Lang: l, Weight: w})
	}

	sort.Sort(acceptLanguages)
	return acceptLanguages
}
