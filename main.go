package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gocolly/colly"
	WebServe "github.com/training_ex1/Training_Program"
)

var port, dir string
var testNum int
var mux sync.Mutex

func init() {
	// flag, set port and dir, one more needed
	flag.StringVar(&port, "port", "8000", "input port")
	flag.StringVar(&dir, "dir", "/Users/leungyantung/go/src/github.com/training_ex1", "input path")

}

func webClamb(file os.FileInfo, coll colly.Collector, coll2 colly.Collector, port string, wg *sync.WaitGroup) {
	if file.IsDir() {
		return // continue
	} else {
		// Print out all file name
		// fmt.Println(file.Name())

		// find "file.name().csv", if not existed, create one; if existed, clear it.
		// the output file will created in path ./output/

		output, err := os.Create("./output/" + strings.Replace(file.Name(), ".html", "", -1) + ".csv")

		//fmt.Println(output.Name())
		if err != nil {
			fmt.Println("File is failed, err: ", err)
		}
		defer output.Close()

		// find class=text-light from <tr>
		coll.OnHTML("table.matrix ", func(eTable *colly.HTMLElement) {
			mux.Lock()
			testNum = 0
			trTableEmpty := true

			eTable.ForEach(".staggered-odd", func(_ int, e *colly.HTMLElement) {
				output.WriteString(";")
				// output history, error least one can't print out
				e.ForEach("div", func(_ int, el *colly.HTMLElement) {
					output.WriteString(el.ChildText("div:nth-child(1)"))
					if el.ChildText("div:nth-child(1)") != "" {
						testNum++
					}
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
				output.WriteString("\n")
			}

			// output date
			eTable.ForEach(".staggered-odd", func(_ int, e *colly.HTMLElement) {
				output.WriteString(";" + e.ChildText("div:nth-child(2)"))
				trTableEmpty = false
			})

			//mux.Unlock()

			output.WriteString("\n;;")

			// output failed info
			output.WriteString("\nfailed;")
			for i := 0; i < testNum; i++ {
				if i < 25 {
					output.WriteString("=COUNTIF(" + string(66+i) + "8:" + string(66+i) + "330, \"x\");")
				} else {
					output.WriteString("=COUNTIF(" + string(65) + string(65+i%25) + "8:" + string(65) + string(65+i%25) + "330, \"x\")")
				}
			}
			//output.WriteString("\nfailed;" + "COUNTIF(B8:B330, \"x\");")
			// output pass info
			output.WriteString("\npass;")
			for i := 0; i < testNum; i++ {
				if i < 25 {
					output.WriteString("=SUM(" + string(66+i) + "8:" + string(66+i) + "330);")
				} else {
					output.WriteString("=SUM(" + string(65) + string(65+i%25) + "8:" + string(65) + string(65+i%25) + "330);")
				}
			}
			// output total run info
			output.WriteString("\ntotal run;")
			for i := 0; i < testNum; i++ {
				if i < 25 {
					output.WriteString("=COUNTA(" + string(66+i) + "8:" + string(66+i) + "330);")
				} else {
					output.WriteString("=COUNTA(" + string(65) + string(65+i%25) + "8:" + string(65) + string(65+i%25) + "330);")
				}
			}
			// output pass rate info
			output.WriteString("\npass rate;")
			for i := 0; i < testNum; i++ {
				if i < 25 {
					output.WriteString("=IF(" + string(66+i) + "5=0, \"N/A\"," + string(66+i) + "5/" + string(66+i) + "6);")
				} else {
					output.WriteString("=IF(" + string(65) + string(65+i%25) + "5=0, \"N/A\"," + string(65) + string(65+i%25) + "5/" + string(65) + string(65+i%25) + "6);")
				}
			}

			output.WriteString("\n;;")

			eTable.ForEach("tr", func(_ int, e *colly.HTMLElement) {

				//trEmpty := true

				//mux.Lock()
				//testNum = 25

				// output model name
				e.ForEach(".boardmodel", func(_ int, el *colly.HTMLElement) {
					if el.ChildText("span") != "" {
						output.WriteString("\n" + el.ChildText("span") + ";")
					}

					// el.ForEach("span", func(_ int, el2 *colly.HTMLElement) {
					// 	output.WriteString(el2.Text + ";")
					// 	trEmpty = false
					// })
				})
				//fmt.Println("in onhtml", file.Name())
				//mux.Unlock()
				//mux.Lock()

				e.ForEach(".cell-full", func(_ int, el *colly.HTMLElement) {
					el.ForEach(".text-success", func(_ int, _ *colly.HTMLElement) {
						output.WriteString("1")
					})
					el.ForEach(".text-light", func(_ int, _ *colly.HTMLElement) {
						output.WriteString("x")
					})
					output.WriteString(";")
				})
				// output to the file with the same name

			})
			mux.Unlock()
		})

		coll.Visit("http://localhost:" + port + "/" + file.Name())
		// fmt.Println(output.Name())
		wg.Done()

	}

}

func main() {
	coll := colly.NewCollector()
	coll2 := colly.NewCollector()
	wg := new(sync.WaitGroup)

	flag.Parse()

	go WebServe.WebServer(port)

	// flag setting
	//flag.StringVar(&port, "port", "")

	data := `./Training_Program/static` // test data path
	files, err := ioutil.ReadDir(data)  // read all the file inside the path above
	if err != nil {
		fmt.Println(err)
	}

	// create a folder called output if it doesnt exist
	outputPath := filepath.Join(".", "output")
	err = os.MkdirAll(outputPath, os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}

	for _, file := range files {
		wg.Add(1)
		go webClamb(file, *coll, *coll2, port, wg)
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
