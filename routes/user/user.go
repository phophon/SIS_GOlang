package user

import (
	"net/http"

	"app"
	"templates"
)

func UserHandler(w http.ResponseWriter, r *http.Request) {
	url := "http://localhost:3010/api"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("authorization", "Bearer YOUR_ACCESS_TOKEN")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))


	// session, err := app.Store.Get(r, "auth-session")
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// templates.RenderTemplate(w, "user", session.Values["profile"])
}
