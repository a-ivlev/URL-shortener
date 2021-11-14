package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

var ts = template.Must(template.New("statPage.html").Funcs(funcMap).ParseFiles(
	"web/statPage.html",
	"web/head.html",
	"web/footer.html",
))

type StatLink struct {
	StatLink string `json:"stat_link"`
}

type Statistic struct {
	Title      string
	ShortLink  string      `json:"short_link"`
	TotalCount int         `json:"total_count"`
	CreatedAt  time.Time   `json:"created_at"`
	FollowList []Following `json:"follow_list"`
	Error      string
}

type Following struct {
	ID              uuid.UUID `json:"id"`
	ShortenerID     uuid.UUID `json:"shortener_id"`
	StatLink        string    `json:"stat_link"`
	IPaddress       string    `json:"ip_address"`
	Count           int       `json:"count"`
	FollowLinkAt time.Time `json:"follow_link_at"`
}

func StatPage(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Error panic: %s (%T)\n", err, err)
		}
	}()
	p := &Statistic{
		Title: "Shortener",
	}

	if r.Method == http.MethodGet {

		var b bytes.Buffer

		statLink := &StatLink{
			// StatLink: r.URL.String(),
			StatLink: chi.URLParam(r, "stat"),
		}

		strJSON, err := json.Marshal(&statLink)
		if err != nil {
			fmt.Fprintf(w, "func StatPage: error occured json marshaling stat page")
		}

		// TODO Хочу эту часть кода перенести в main но незнаю как добавить контекст,
		// чтобы здесь это получить через r.Context().Value("SRV_HOST")
		srvHost := os.Getenv("SRV_HOST")
		if srvHost == "" {
			log.Fatal("unknown SRV_HOST = ", srvHost)
		}
		//srvPort := os.Getenv("SRV_PORT")
		//if srvPort == "" {
		//	log.Fatal("unknown SRV_PORT = ", srvPort)
		//}
		//-------------------------------------------------------------------------------

		//srv := fmt.Sprintf("http://%s:%s/stat", srvHost, srvPort)
		srv := fmt.Sprintf("http://%s/stat", srvHost)

		client := &http.Client{Timeout: time.Second * 5}
		req, err := http.NewRequest(http.MethodPost, srv, bytes.NewBuffer(strJSON))
		if err != nil {
			log.Println("func StatPage: error occurred NewRequest: ", err)
		}
		req.Header.Set("Content-Type", "application/json")

		res, err := client.Do(req)
		if err != nil {
			log.Println("func StatPage: error occurred client Do: ", err)
			p.Error = "func StatPage: error occurred client Do"
			b.Reset()
			err = ts.ExecuteTemplate(&b, "statPage.html", p)
			if err != nil {
				fmt.Fprintf(w, "func StatPage: a error occured rendering template statPage.html")
				log.Println("func StatPage: a error occurred rendering template statPage.html: ", err)
				return
			}
			_, err = b.WriteTo(w)
			if err != nil {
				log.Println("func StatPage: writing error home page: ", err)
			}
		}

		statDB := &Statistic{}
		defer res.Body.Close()
		if err = json.NewDecoder(res.Body).Decode(&statDB); err != nil {
			http.Error(w, "error unmarshal request", http.StatusInternalServerError)
			return
		}

		// TODO Хочу эту часть кода перенести в main но незнаю как добавить контекст,
		// чтобы здесь это получить через r.Context().Value("SHORT_CLI_HOST")
		cliHost := os.Getenv("CLI_HOST")
		if cliHost == "" {
			log.Fatal("unknown CLI_HOST = ", cliHost)
		}
		//cliPort := os.Getenv("CLI_PORT")
		//if cliPort == "" {
		//	log.Fatal("unknown PORT = ", cliPort)
		//}

		//p.ShortLink = fmt.Sprintf("http://%s:%s/%s", cliHost, cliPort, statDB.ShortLink)
		p.ShortLink = fmt.Sprintf("http://%s/%s", cliHost, statDB.ShortLink)
		p.TotalCount = statDB.TotalCount
		p.CreatedAt = statDB.CreatedAt
		p.FollowList = statDB.FollowList

		err = ts.ExecuteTemplate(&b, "statPage.html", p)
		if err != nil {
			fmt.Fprintf(w, "A error occured execute template statPage.html .")
			return
		}
		_, err = b.WriteTo(w)
		if err != nil {
			log.Println("write render home page error: ", err)
		}
	}
}
