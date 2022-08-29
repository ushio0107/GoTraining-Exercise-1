package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/gocolly/colly"
	WebServe "github.com/training_ex1/Training_Program"
)

type Info struct {
	date      string
	modelName string
	modelFail []int
}

func webClamb(mux *sync.Mutex, file os.FileInfo, coll colly.Collector, port string, wg *sync.WaitGroup) {
	//for _, file := range files {
	if file.IsDir() {
		return // continue
	} else {
		fmt.Println(file.Name())

		// find "file.name().csv", if not existed, create one; if existed, clear it.
		// the output file will created in path ./output/

		output, err := os.Create("./output/" + strings.Replace(file.Name(), ".html", "", -1) + ".csv")

		//fmt.Println(output.Name())
		if err != nil {
			fmt.Println("File is failed, err: ", err)
		}
		output.WriteString(";")

		// find class=text-light from <tr>
		coll.OnHTML("table", func(eTable *colly.HTMLElement) {
			trTableEmpty := true
			mux.Lock()
			eTable.ForEach(".staggered-odd", func(_ int, e *colly.HTMLElement) {
				output.WriteString(";")
				// output history, error least one can't print out
				e.ForEach("div", func(_ int, el *colly.HTMLElement) {
					output.WriteString(el.ChildText("div:nth-child(1)"))
					trTableEmpty = false
				})
				//output.WriteString(";" + e.ChildText("div:nth-child(1)"))

				// e.ForEach("a[href]", func(_ int, el *colly.HTMLElement) {
				// 	output.WriteString(";" + el.Text)
				// 	trTableEmpty = false
				// })

				// output.WriteString(" : " + e.ChildText("div:nth-child(2)"))

				// e.ForEach("div > div > div", func(_ int, el *colly.HTMLElement) {
				// 	output.WriteString(" : " + el.Text)
				// 	trTableEmpty = false
				// })
			})

			if trTableEmpty == false {
				output.WriteString("\n;")
			}

			// output date
			eTable.ForEach(".staggered-odd", func(_ int, e *colly.HTMLElement) {
				output.WriteString(";" + e.ChildText("div:nth-child(2)"))
				trTableEmpty = false
			})

			if trTableEmpty == false {
				output.WriteString("\n")
			}
			mux.Unlock()

			eTable.ForEach("tr", func(_ int, e *colly.HTMLElement) {
				trEmpty := true
				mux.Lock()

				// output model name
				e.ForEach(".boardmodel", func(_ int, el *colly.HTMLElement) {
					el.ForEach("span", func(_ int, el2 *colly.HTMLElement) {
						output.WriteString(el2.Text + ";")
						trEmpty = false
					})
				})
				//fmt.Println("in onhtml", file.Name())
				e.ForEach(".cell-full", func(_ int, el *colly.HTMLElement) {
					el.ForEach(".text-light", func(_ int, el2 *colly.HTMLElement) {
						output.WriteString(el2.Text)
					})
					output.WriteString(";")
				})
				if trEmpty == false {
					output.WriteString("\n")
				}
				// output to the file with the same name
				mux.Unlock()
			})

		})

		coll.Visit("http://localhost:" + port + "/" + file.Name())
		// fmt.Println(output.Name())
		output.Close()
		wg.Done()
	}

	//}
}

func main() {
	coll := colly.NewCollector()
	var mux sync.Mutex
	wg := new(sync.WaitGroup)
	wg.Add(68)

	//var port string
	var port, dir string
	flag.StringVar(&port, "port", "8000", "input port")
	flag.StringVar(&dir, "dir", "/Users/leungyantung/go/src/github.com/training_ex1", "input path")
	flag.Parse()

	go WebServe.WebServer(port)

	// flag setting
	//flag.StringVar(&port, "port", "")

	data := `./Training_Program/static` // test data path
	files, err := ioutil.ReadDir(data)  // read all the file inside the path above
	if err != nil {
		fmt.Println(err)
	}

	for _, file := range files {
		go webClamb(&mux, file, *coll, port, wg)
	}

	wg.Wait()

	// find "name.csv", if not existed, create one; if existed, clear it.
	// output, err := os.Create("arc.MediaSourceUI" + ".csv")
	// if err != nil {
	// 	fmt.Println("File is failed, err: ", err)
	// }
	// defer output.Close()
	// go func() {
	// 	for _, file := range files {
	// 		if file.IsDir() {
	// 			continue
	// 		} else {
	// 			fmt.Println(file.Name(), i)
	// 			i = i + 1

	// 			// find "file.name().csv", if not existed, create one; if existed, clear it.
	// 			// the output file will created in path ./output/

	// 			output, err := os.Create("./output/" + strings.Replace(file.Name(), ".html", "", -1) + ".csv")

	// 			//fmt.Println(output.Name())
	// 			if err != nil {
	// 				fmt.Println("File is failed, err: ", err)
	// 			}
	// 			output.WriteString(";")

	// 			// find class=text-light from <tr>
	// 			coll.OnHTML("table", func(eTable *colly.HTMLElement) {
	// 				trTableEmpty := true
	// 				mux.Lock()
	// 				eTable.ForEach(".staggered-odd", func(_ int, e *colly.HTMLElement) {
	// 					// output history, error least one can't print out
	// 					e.ForEach("a[href]", func(_ int, el *colly.HTMLElement) {
	// 						output.WriteString(";123" + el.Text)
	// 						trTableEmpty = false
	// 					})

	// 					// output.WriteString(" : " + e.ChildText("div:nth-child(2)"))

	// 					// e.ForEach("div > div > div", func(_ int, el *colly.HTMLElement) {
	// 					// 	output.WriteString(" : " + el.Text)
	// 					// 	trTableEmpty = false
	// 					// })
	// 				})

	// 				if trTableEmpty == false {
	// 					output.WriteString("\n;")
	// 				}

	// 				// output date
	// 				eTable.ForEach(".staggered-odd", func(_ int, e *colly.HTMLElement) {
	// 					output.WriteString(";" + e.ChildText("div:nth-child(2)"))
	// 					trTableEmpty = false
	// 				})

	// 				if trTableEmpty == false {
	// 					output.WriteString("\n")
	// 				}
	// 				mux.Unlock()

	// 				eTable.ForEach("tr", func(_ int, e *colly.HTMLElement) {
	// 					trEmpty := true
	// 					mux.Lock()

	// 					// output model name
	// 					e.ForEach(".boardmodel", func(_ int, el *colly.HTMLElement) {
	// 						el.ForEach("span", func(_ int, el2 *colly.HTMLElement) {
	// 							output.WriteString(el2.Text + ";")
	// 							trEmpty = false
	// 						})
	// 					})
	// 					//fmt.Println("in onhtml", file.Name())
	// 					e.ForEach(".cell-full", func(_ int, el *colly.HTMLElement) {
	// 						el.ForEach(".text-light", func(_ int, el2 *colly.HTMLElement) {
	// 							output.WriteString(el2.Text)
	// 						})
	// 						output.WriteString(";")
	// 					})
	// 					if trEmpty == false {
	// 						output.WriteString("\n")
	// 					}
	// 					// output to the file with the same name
	// 					mux.Unlock()
	// 				})

	// 			})

	// 			coll.Visit("http://localhost:" + port + "/" + file.Name())
	// 			// fmt.Println(output.Name())
	// 			output.Close()

	// 		}

	// 	}
	// }()

	// testWrite, err := os.Create("./output/" + "arc.MediaSourceUI.html" + ".csv")

	// // find class=text-light from <tr>
	// coll.OnHTML("tr", func(e *colly.HTMLElement) {
	// 	testWrite.WriteString("module;")
	// 	e.ForEach(".text-light", func(_ int, el *colly.HTMLElement) {
	// 		if el.Text != "" {
	// 			fmt.Print(el.Text, ";")
	// 			testWrite.WriteString(el.Text + ";")
	// 		}
	// 	})
	// 	fmt.Println(";")
	// 	// output to the file with the same name
	// 	testWrite.WriteString(";\n")
	// })
	//time.Sleep(25 * time.Second)

	// coll.Visit("http://localhost:" + port + "/" + "arc.MediaSourceUI.html")

}
