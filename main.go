package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/gocolly/colly"
	WebServe "github.com/training_ex1/Training_Program"
)

var port, dir, outputDirName string
var testNum, modelNum int
var mux sync.Mutex

func init() {
	// flag, set port and dir, one more needed
	flag.StringVar(&port, "port", "8000", "input port")
	flag.StringVar(&dir, "dir", "/Users/leungyantung/go/src/github.com/training_ex1", "input path")
	flag.StringVar(&outputDirName, "outputDirName", "/output", "input outputDirName")

}

func webClamb(outputPath string, file os.FileInfo, collUpper colly.Collector, collLower colly.Collector, port string, wg *sync.WaitGroup) {
	if file.IsDir() {
		return // continue
	} else {
		// Print out all file name
		// fmt.Println(file.Name())

		// find "file.name().csv", if not existed, create one; if existed, clear it.
		// the output file will created in path ./output/

		output, err := os.Create(outputPath + "/" + strings.Replace(file.Name(), ".html", "", -1) + ".csv")

		//fmt.Println(output.Name())
		if err != nil {
			fmt.Println("File is failed, err: ", err)
		}
		defer output.Close()

		// find class=text-light from <tr>
		collUpper.OnHTML("table.matrix ", func(eTable *colly.HTMLElement) {
			mux.Lock()
			testNum, modelNum = 0, 0
			trTableEmpty := true

			eTable.ForEach("tr", func(_ int, e *colly.HTMLElement) {
				e.ForEach(".boardmodel", func(_ int, el *colly.HTMLElement) {
					if el.ChildText("span") != "" {
						modelNum++
					}
				})
			})

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
					output.WriteString("=COUNTIF(" + string(66+i) + "9:" + string(66+i) + strconv.Itoa(modelNum+8) + ", \"x\");")
				} else {
					output.WriteString("=COUNTIF(" + string(65) + string(65+i%25) + "9:" + string(65) + string(65+i%25) + strconv.Itoa(modelNum+8) + ", \"x\");")
				}
			}
			//output.WriteString("\nfailed;" + "COUNTIF(B9:B330, \"x\");")
			// output pass info
			output.WriteString("\npass;")
			for i := 0; i < testNum; i++ {
				if i < 25 {
					output.WriteString("=SUM(" + string(66+i) + "9:" + string(66+i) + strconv.Itoa(modelNum+8) + ");")
				} else {
					output.WriteString("=SUM(" + string(65) + string(65+i%25) + "9:" + string(65) + string(65+i%25) + strconv.Itoa(modelNum+8) + ");")
				}
			}
			// output total run info
			output.WriteString("\ntotal run;")
			for i := 0; i < testNum; i++ {
				if i < 25 {
					output.WriteString("=COUNTA(" + string(66+i) + "9:" + string(66+i) + strconv.Itoa(modelNum+8) + ");")
				} else {
					output.WriteString("=COUNTA(" + string(65) + string(65+i%25) + "9:" + string(65) + string(65+i%25) + strconv.Itoa(modelNum+8) + ");")
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
		})

		collUpper.Visit("http://localhost:" + port + "/" + file.Name())

		collLower.OnHTML("table.matrix ", func(eTable *colly.HTMLElement) {
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

				e.ForEach(".cell-full", func(_ int, el *colly.HTMLElement) {
					el.ForEach(".text-success", func(_ int, _ *colly.HTMLElement) {
						output.WriteString("1")
					})
					el.ForEach(".text-light", func(_ int, _ *colly.HTMLElement) {
						output.WriteString("x")
					})
					output.WriteString(";")
				})

			})
			mux.Unlock()
		})

		collLower.Visit("http://localhost:" + port + "/" + file.Name())
		// fmt.Println(output.Name())
		wg.Done()

	}

}

func main() {
	collUpper := colly.NewCollector()
	collLower := colly.NewCollector()
	wg := new(sync.WaitGroup)

	flag.Parse()

	data := `/test_data`                     // test data path
	files, err := ioutil.ReadDir(dir + data) // read all the file inside the path above
	if err != nil {
		fmt.Println(err)
	}
	go WebServe.WebServer(port, dir+data)

	// create a folder called output if it doesnt exist
	outputPath := filepath.Join(dir, outputDirName)
	err = os.MkdirAll(outputPath, os.ModePerm)
	fmt.Println(outputPath)
	if err != nil {
		fmt.Println(err)
	}

	for _, file := range files {
		wg.Add(1)
		go webClamb(outputPath, file, *collUpper, *collLower, port, wg)
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
