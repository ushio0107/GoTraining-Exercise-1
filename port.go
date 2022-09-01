// package main

// import (
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"net/http/httptest"
// 	"os"
// 	"time"

// 	"github.com/gocolly/colly"
// )

// // func main() {
// // 	fmt.Println("inside Webserver")
// // 	//var port string

// // 	req, err := http.NewRequest("GET", "http://localhost:8000", nil)

// // 	// flag.StringVar(&port, "port", "8000", "input port")
// // 	// flag.Parse()
// // 	// fmt.Println(port)
// // 	// // serve the html file inside ./static
// // 	// http.Handle("/", http.FileServer(http.Dir("../training_ex1/Training_Program/static")))

// // 	// // port 8000
// // 	// err := http.ListenAndServe(":"+port, nil)
// // 	// if err != nil {
// // 	// 	fmt.Println(err)
// // 	// }
// // }

// func main() {
// 	file, err := os.Open("./arc.MediaSourceUI.html")
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	defer file.Close()
// 	// fmt.Println(file.Name())

// 	handle := func(w http.ResponseWriter, r *http.Request) {

// 	}

// 	var r io.Reader
// 	r = file
// 	fmt.Println(r)
// 	req, err2 := http.NewRequest(http.MethodGet, "http://localhost:8000", r)
// 	if err2 != nil {
// 		fmt.Println(err2)
// 	}
// 	rec := httptest.NewRecorder()
// 	http.DefaultServeMux.ServeHTTP(rec, req)

// 	collUpper := colly.NewCollector()

// 	collUpper.OnHTML("table.matrix ", func(eTable *colly.HTMLElement) {
// 		fmt.Println("Inside")
// 	})

// 	collUpper.Visit("http://localhost:8000")
// 	// ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 	// 	fmt.Fprintln(w, "Hello, client")
// 	// }))
// 	// defer ts.Close()

// 	// res, err := http.Get(ts.URL)
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }
// 	// greeting, err := ioutil.ReadAll(res.Body)
// 	// res.Body.Close()
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }

// 	// fmt.Printf("%s", greeting)
// 	time.Sleep(10 * time.Second)
// }

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"

	"github.com/gocolly/colly"
)

var (
	port2          string
	dir2           string
	outputDirName2 string

	mux2 sync.Mutex
)

func init() {
	flag.StringVar(&port2, "port", "8000", "local host server port")
	flag.StringVar(&dir2, "dir", ".", "the working directory ")
	flag.StringVar(&outputDirName2, "outputDirName", "/output", "the name of the directory which stored the output files")

	flag.Parse()

}

func main() {
	// handler := func(w http.ResponseWriter, r *http.Request) {
	// 	file, err := os.Open("./test/arc.MediaSourceUI.html")
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// 	b, err := ioutil.ReadAll(file)
	// 	w.Write(b)
	// }
	files, _ := ioutil.ReadDir("./test_data") // read all the file inside the path above
	fmt.Println(files)
	// if err != nil {
	// 	fmt.Println("1", err)
	// }
	for _, file := range files {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			file2, err := os.Open(dir2 + "/test_data/" + file.Name())
			if err != nil {
				fmt.Println("2", err)
			}
			b, err := ioutil.ReadAll(file2)
			w.Write(b)
		}))
		collUpper := colly.NewCollector()

		collUpper.OnHTML("table.matrix ", func(eTable *colly.HTMLElement) {
			fmt.Println("Inside")
		})

		collUpper.Visit(ts.URL)

		// u, err := url.Parse(ts.URL)
		// if err != nil {
		// 	fmt.Println("3", err)
		// }
		// fmt.Println(u)
	}

	// collUpper := colly.NewCollector()

	// collUpper.OnHTML("table.matrix ", func(eTable *colly.HTMLElement) {
	// 	//fmt.Println("Inside")
	// })

	// collUpper.Visit("http://localhost:8000")

	// log.Println(string(body))
	//fmt.Println(string(body))

}
