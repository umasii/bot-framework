package walmart

import (
	"fmt"
	"math/rand"
	"time"
)

var acceptList = []string{
	",application/xhtml+xml",
	",application/xml",
	",image/webp",
	",image/apng",
	",text/html",
}

func randomAccept() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf(`application/json%s`, acceptList[rand.Intn(len(acceptList))])
}

func randomChUa() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf(`" Not A;Brand";v="%d", "Chromium";v="%d", "Google Chrome";v="%d"`, rand.Intn(10)+90, rand.Intn(12)+80, rand.Intn(12)+80)
}

var langList = []string{
	"fr-CH",
	"fr",
	"de",
	"*",
	"en-GB",
	"de-CH",
	"it",
	"es",
}

func randomAcceptLanguage() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("en-US,en;q=0.9,%s;q=0.%d", langList[rand.Intn(len(langList))], rand.Intn(9)+1)
}

func randomizeHeaders(headers []map[string]string) []map[string]string {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(headers), func(i, j int) {
		headers[i], headers[j] = headers[j], headers[i]
	})

	return headers

}
