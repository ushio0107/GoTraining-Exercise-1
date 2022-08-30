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

type Info struct {
	testCase    string
	testDate    string
	testNum     int
	modelNum    int
	modelName   string
	excelFomula string
	testResult  string
}

func init() {
	// flag, set port and dir, one more needed
	flag.StringVar(&port, "port", "8000", "input port")
	flag.StringVar(&dir, "dir", "/Users/leungyantung/go/src/github.com/training_ex1", "input path")
	flag.StringVar(&outputDirName, "outputDirName", "/output", "input outputDirName")

}

func webClamb(testCase Info, outputPath string, file os.FileInfo, collUpper colly.Collector, collLower colly.Collector, port string, wg *sync.WaitGroup) {
	if file.IsDir() {
		return // continue
	} else {

		// To avoid output files aren't exist, os.Create will create file if file doesn't exist,
		// clear file if file already existed, output file name will as the same as input file
		// also, to avoid double file format, .html will be taken out by strings.Replace
		output, err := os.Create(outputPath + "/" + strings.Replace(file.Name(), ".html", "", -1) + ".csv")
		fmt.Println("Created output file: ", output.Name())
		if err != nil {
			fmt.Println("File is failed, err: ", err)
		}
		defer output.Close()

		// select path
		collUpper.OnHTML("table.matrix ", func(eTable *colly.HTMLElement) {
			mux.Lock()
			testCase.modelNum, testCase.testNum = 0, 0
			//testNum, modelNum = 0, 0
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
					if el.ChildText("div:nth-child(1)") != "" {
						output.WriteString(" " + el.ChildText("div:nth-child(1)") + " ")
						testNum++
					}
					trTableEmpty = false
				})
			})

			if trTableEmpty == false {
				output.WriteString("\n")
			}

			// output date
			eTable.ForEach(".staggered-odd", func(_ int, e *colly.HTMLElement) {
				output.WriteString(";" + " " + e.ChildText("div:nth-child(2)") + " ")
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
			mux.Unlock()
		})

		collUpper.Visit("http://localhost:" + port + "/" + file.Name())

		fmt.Println(file.Name(), " ", modelNum)

		collLower.OnHTML("table.matrix ", func(eTable *colly.HTMLElement) {
			eTable.ForEach("tr", func(_ int, e *colly.HTMLElement) {

				// output model name
				e.ForEach(".boardmodel", func(_ int, el *colly.HTMLElement) {
					if el.ChildText("span") != "" {
						output.WriteString("\n" + el.ChildText("span") + ";")
					}

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
		})

		collLower.Visit("http://localhost:" + port + "/" + file.Name())
		//fmt.Println(output.Name())
		wg.Done()

	}

}

func main() {
	defer os.Exit(0)

	var testCase Info

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
	if err != nil {
		fmt.Println(err)
	}

	for _, file := range files {
		wg.Add(1)
		go webClamb(testCase, outputPath, file, *collUpper, *collLower, port, wg)
	}

	wg.Wait()

}
