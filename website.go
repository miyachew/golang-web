package main

import ("fmt"
				"net/http"
				"html/template"
				"io/ioutil"
				"encoding/xml"				
)

type Sitemapindex struct {
	Locations []string `xml:"sitemap>loc"`
}

type News struct {
	Titles []string `xml:"url>news>title"`
	Keywords []string `xml:"url>news>keywords"`
	Locations []string `xml:"url>loc"`
}

type FormattedNews struct {
	Keyword string
	Location string
}

type NewsAggPage struct {
    Title string
    News map[string]FormattedNews
}

func indexHandler(w http.ResponseWriter,r *http.Request){
	fmt.Fprintf(w, `go is neat!
	sddsfsf`)
	fmt.Fprintf(w,"heeeeello %s","miya")
}

const newsTemplate string = "news.html"

func newsHandler(w http.ResponseWriter,r *http.Request){
	newsMap := getNewsMap("https://www.washingtonpost.com/news-sitemap-index.xml")
	p := NewsAggPage{ Title : "News from washingtonpost.com", News : newsMap}
	tmpl, _ := template.ParseFiles(newsTemplate)
	tmpl.Execute(w,p)
}

func getUrlBytes(url string) []byte{
	res, _ := http.Get(url)
	bytes, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	return bytes
}

func getNewsMap(url string) map[string]FormattedNews{
	MainXmlBytes := getUrlBytes(url)

	var s Sitemapindex
	xml.Unmarshal(MainXmlBytes, &s)

	var news News
	newsMap := make(map[string]FormattedNews)

	for _, location := range s.Locations {
		bytes := getUrlBytes(location)
		xml.Unmarshal(bytes,&news)
		
		for k, v := range news.Titles {
				newsMap[v] = FormattedNews{ news.Keywords[k],news.Locations[k] }
		}
		break
	}
	return newsMap
}

func main(){
	http.HandleFunc("/",indexHandler)
	http.HandleFunc("/news",newsHandler)
	http.ListenAndServe(":8000",nil)
}