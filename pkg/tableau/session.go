package tableau

import (
	"fmt"
	stdhttp "net/http"

	"github.com/AdhityaRamadhanus/cockpit/pkg/http"
)

func GetSessionID(workBook string, dashboard string) (string, error) {
	// url := "https://public.tableau.com/views/DashboardCovid-19Jakarta_15837354399300/Dashboard22?%3Aembed=y&%3AshowVizHome=no&%3Adisplay_count=y&%3Adisplay_static_image=y&%3AbootstrapWhenNotified=true"
	// url := "https://public.tableau.com/views/PetaPersebaranTes/Dashboard2?%3Aembed=y&%3AshowVizHome=no&%3Adisplay_count=y&%3Adisplay_static_image=y&%3AbootstrapWhenNotified=true"
	url := fmt.Sprintf(
		"https://public.tableau.com/views/%s/%s?",
		workBook,
		dashboard,
	)
	url += "%3Aembed=y&%3AshowVizHome=no&%3Adisplay_count=y&%3Adisplay_static_image=y&%3AbootstrapWhenNotified=true"

	req, err := stdhttp.NewRequest(stdhttp.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("authority", "public.tableau.com")
	req.Header.Set("upgrade-insecure-requests", "1")
	req.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.79 Safari/537.36")
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sec-fetch-mode", "nested-navigate")
	req.Header.Set("referer", "/public.tableau.com/profile/jsc.data")
	req.Header.Set("accept-encoding", "gzip, deflate, br")
	req.Header.Set("accept-language", "en-US,en;q=0.9")
	req.Header.Set("cookie", "tableau_locale=en; tableau_public_negotiated_locale=en-us; _ga=GA1.2.2001330615.1584503744; _gid=GA1.2.67936358.1584503744; _ga=GA1.3.2001330615.1584503744; _gid=GA1.3.67936358.1584503744; has_js=1; ELOQUA=GUID=4733810149F8461D8F39004E928855C0; _gd_svisitor=37fb3d1760640000b1923a5e8f030000bd070000; _gcl_au=1.1.2052117600.1584505310; _fbp=fb.1.1584505310896.778815216; _dc_gtm_UA-625217-47=1; _gat_UA-625217-22=1; _gat_UA-625217-47=1")

	client := http.NewDefaultRetryClient()
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	return res.Header.Get("x-session-id"), nil
}
