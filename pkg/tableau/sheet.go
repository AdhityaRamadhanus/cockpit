package tableau

import (
	"io/ioutil"
	stdhttp "net/http"
	stdurl "net/url"
	"strings"

	"github.com/AdhityaRamadhanus/cockpit/pkg/http"
)

func GetSheet(url string, sheetID string) (string, error) {
	form := stdurl.Values{}
	form.Add("sheet_id", sheetID)

	req, err := stdhttp.NewRequest(stdhttp.MethodPost, url, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("authority", "public.tableau.com")
	req.Header.Set("accept", "text/plain")
	req.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.79 Safari/537.36")
	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	req.Header.Set("origin", "/public.tableau.com")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("accept-language", "en-US,en;q=0.9")
	req.Header.Set("cookie", "tableau_locale=en; tableau_public_negotiated_locale=en-us; _ga=GA1.2.2001330615.1584503744; _gid=GA1.2.67936358.1584503744; _ga=GA1.3.2001330615.1584503744; _gid=GA1.3.67936358.1584503744; has_js=1; ELOQUA=GUID=4733810149F8461D8F39004E928855C0; _gd_svisitor=37fb3d1760640000b1923a5e8f030000bd070000; _gcl_au=1.1.2052117600.1584505310; _fbp=fb.1.1584505310896.778815216; _gat_UA-625217-22=1; _gat_UA-625217-47=1; _dc_gtm_UA-625217-47=1")

	client := http.NewDefaultRetryClient()
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	return string(body), nil
}
