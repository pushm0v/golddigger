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

	// Coba navigasi ke elemen div yang berisi IDR/g berdasar struktur Next.js terbaru
	doc.Find("div").Each(func(i int, s *goquery.Selection) {
		if price > 0 {
			return // Sudah ketemu nilainya
		}

		// Cari elemen span yang menampilkan text "IDR/g"
		if s.Find("span").First().Text() == "IDR/g" {
			rawText := s.Text()

			// Memisahkan teks untuk mendapatkan value-nya
			// Bentuk string dari raw HTML biasanya: "IDR/gRp2.867.825+Rp0 (+0.00%)"
			if strings.Contains(rawText, "IDR/gRp") {
				cutString := strings.Split(rawText, "IDR/gRp")
				if len(cutString) > 1 {
					valString := cutString[1]

					// Potong string yang ada tanda '+' (kenaikan harga emas harian)
					cutKenaikan := strings.Split(valString, "+")
					if len(cutKenaikan) > 0 {
						// Ekstrak angka murni (buang pemisah titik)
						reg, regexErr := regexp.Compile("[^0-9]+")
						if regexErr != nil {
							err = regexErr
							return
						}

						processedString := reg.ReplaceAllString(cutKenaikan[0], "")
						p, parseErr := strconv.Atoi(processedString)
						if parseErr == nil && p > 0 {
							price = p
						}
					}
				}
			}
		}
	})

	// Fallback error apabila tidak ketemu sama sekali (kemungkinan struktur HTML ganti lagi)
	if price == 0 && err == nil {
		err = fmt.Errorf("failed to parse gold price from the retrieved html")
	}

	return
}
