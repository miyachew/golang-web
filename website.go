package main

import ("fmt"
				"net/http"
				"html/template"
				"io/ioutil"
				"encoding/xml"
				"sync"
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
var wg sync.WaitGroup

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

func newsRoutine(c chan News, url string){
	defer wg.Done()
	var news News
	res, _ := http.Get(url)
	bytes, _ := ioutil.ReadAll(res.Body)
	xml.Unmarshal(bytes,&news)
	res.Body.Close()

	c <- news
}

func getNewsMap(url string) map[string]FormattedNews{
	MainXmlBytes := getUrlBytes(url)
	var s Sitemapindex
	xml.Unmarshal(MainXmlBytes, &s)

	queue := make(chan News, 30)
	
	newsMap := make(map[string]FormattedNews)
	
	for _, Location := range s.Locations {
		wg.Add(1)
		go newsRoutine(queue, Location)
	}
	wg.Wait()
	close(queue)

	for elem := range queue {
		for k, v := range elem.Titles {
				newsMap[v] = FormattedNews{ elem.Keywords[k],elem.Locations[k] }
		}
	}
	return newsMap
}

func main(){
	http.HandleFunc("/",indexHandler)
	http.HandleFunc("/news",newsHandler)
	http.ListenAndServe(":8000",nil)
}