package golddigger

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func DigHargaEmasOrg() (price int, err error) {
	price = 0
	// Request the HTML page.
	res, err := http.Get("https://harga-emas.org/1-gram/")
	if err != nil {
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		err = fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
		return
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return
	}

	doc.Find(".in_table tr").Each(func(i int, s *goquery.Selection) {
		if i == 4 {
			s.Find("td").Each(func(j int, t *goquery.Selection) {
				if j == 1 {
					rawText := t.Text()
					reg, err := regexp.Compile("[^0-9,]+")
					if err != nil {
						return
					}
					processedString := reg.ReplaceAllString(rawText, "")
					parts := strings.Split(processedString, ",")
					if len(parts) >= 1 {
						price, err = strconv.Atoi(parts[0])
						if err != nil {
							return
						}
					}
				}
			})
		}
	})
	return
}
