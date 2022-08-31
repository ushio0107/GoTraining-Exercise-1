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
	testInfo    TestInfo
	modelInfo   ModelInfo
	excelFomula string
}

type TestInfo struct {
	testNo     []string
	testDate   []string
	testNum    int
	testResult []string
}

type ModelInfo struct {
	modelNum  int
	modelName []string
}

func init() {
	// flag, set port and dir, one more needed
	flag.StringVar(&port, "port", "8000", "input port")
	flag.StringVar(&dir, "dir", ".", "input path")
	flag.StringVar(&outputDirName, "outputDirName", "/output", "input outputDirName")

}

func webClamb(testCase *Info, outputPath string, file os.FileInfo, collUpper colly.Collector, collLower colly.Collector, port string, wg *sync.WaitGroup) {
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
			testNum, modelNum = 0, 0
			trTableEmpty := true

			//body > div > table.matrix.table.table-sm.table-less-padding.table-borderless.table-striped > tbody > tr:nth-child(3) > td.boardmodel.text-ellipsis > span
			// eTable.ForEach("tr", func(_ int, e *colly.HTMLElement) {
			// 	e.ForEach(".boardmodel", func(_ int, el *colly.HTMLElement) {
			// 		if el.ChildText("span") != "" {
			// 			modelNum++
			// 		}
			// 	})
			// })
			eTable.ForEach("td.boardmodel", func(_ int, e *colly.HTMLElement) {
				if e.ChildText("span") != "" {
					testCase.modelInfo.modelNum++
					modelNum++
				}
			})

			eTable.ForEach(".staggered-odd", func(_ int, e *colly.HTMLElement) {
				// output.WriteString(";")

				// output history, error least one can't print out
				e.ForEach("div", func(_ int, el *colly.HTMLElement) {
					if el.ChildText("div:nth-child(1)") != "" {
						// output.WriteString(" " + el.ChildText("div:nth-child(1)") + " ")
						//testCase.testNo = el.ChildText("div:nth-child(1)")
						testCase.testInfo.testNo = append(testCase.testInfo.testNo, el.ChildText("div:nth-child(1)"))

						//output.WriteString(" " + testCase.testInfo.testNo[testCase.testInfo.testNum])
						testNum++
						testCase.testInfo.testNum++
					}
					trTableEmpty = false
				})
			})

			if trTableEmpty == false {
				// output.WriteString("\n")
			}

			// output date
			eTable.ForEach(".staggered-odd", func(_ int, e *colly.HTMLElement) {
				testCase.testInfo.testDate = append(testCase.testInfo.testDate, e.ChildText("div:nth-child(2)"))
				// output.WriteString(";" + " " + e.ChildText("div:nth-child(2)") + " ")
				trTableEmpty = false
			})

			//mux.Unlock()

			// output.WriteString("\n;;")

			// output failed info
			// output.WriteString("\nfailed;")
			testCase.excelFomula = "\nfailed;"
			for i := 0; i < testNum; i++ {
				if i < 25 {
					testCase.excelFomula = testCase.excelFomula + "=COUNTIF(" + string(66+i) + "9:" + string(66+i) + strconv.Itoa(modelNum+8) + ", \"x\");"
					// output.WriteString("=COUNTIF(" + string(66+i) + "9:" + string(66+i) + strconv.Itoa(modelNum+8) + ", \"x\");")
				} else {
					testCase.excelFomula = testCase.excelFomula + "=COUNTIF(" + string(65) + string(65+i%25) + "9:" + string(65) + string(65+i%25) + strconv.Itoa(modelNum+8) + ", \"x\");"
					// output.WriteString("=COUNTIF(" + string(65) + string(65+i%25) + "9:" + string(65) + string(65+i%25) + strconv.Itoa(modelNum+8) + ", \"x\");")
				}
			}
			testCase.excelFomula = testCase.excelFomula + "\npass;"
			//output.WriteString("\nfailed;" + "COUNTIF(B9:B330, \"x\");")
			// output pass info
			// output.WriteString("\npass;")
			for i := 0; i < testNum; i++ {
				if i < 25 {
					testCase.excelFomula = testCase.excelFomula + "=SUM(" + string(66+i) + "9:" + string(66+i) + strconv.Itoa(modelNum+8) + ");"
				} else {
					testCase.excelFomula = testCase.excelFomula + "=SUM(" + string(65) + string(65+i%25) + "9:" + string(65) + string(65+i%25) + strconv.Itoa(modelNum+8) + ");"
				}
			}
			// output total run info
			testCase.excelFomula = testCase.excelFomula + "\ntotal run;"
			// output.WriteString("\ntotal run;")
			for i := 0; i < testNum; i++ {
				if i < 25 {
					testCase.excelFomula = testCase.excelFomula + "=COUNTA(" + string(66+i) + "9:" + string(66+i) + strconv.Itoa(modelNum+8) + ");"
					// output.WriteString("=COUNTA(" + string(66+i) + "9:" + string(66+i) + strconv.Itoa(modelNum+8) + ");")
				} else {
					testCase.excelFomula = testCase.excelFomula + "=COUNTA(" + string(65) + string(65+i%25) + "9:" + string(65) + string(65+i%25) + strconv.Itoa(modelNum+8) + ");"
					// output.WriteString("=COUNTA(" + string(65) + string(65+i%25) + "9:" + string(65) + string(65+i%25) + strconv.Itoa(modelNum+8) + ");")
				}
			}
			// output pass rate info
			// output.WriteString("\npass rate;")
			testCase.excelFomula = testCase.excelFomula + "\npass rate;"
			for i := 0; i < testNum; i++ {
				if i < 25 {
					testCase.excelFomula = testCase.excelFomula + "=IF(" + string(66+i) + "5=0, \"N/A\"," + string(66+i) + "5/" + string(66+i) + "6);"
					// output.WriteString("=IF(" + string(66+i) + "5=0, \"N/A\"," + string(66+i) + "5/" + string(66+i) + "6);")
				} else {
					testCase.excelFomula = testCase.excelFomula + "=IF(" + string(65) + string(65+i%25) + "5=0, \"N/A\"," + string(65) + string(65+i%25) + "5/" + string(65) + string(65+i%25) + "6);"
					// output.WriteString("=IF(" + string(65) + string(65+i%25) + "5=0, \"N/A\"," + string(65) + string(65+i%25) + "5/" + string(65) + string(65+i%25) + "6);")
				}
			}

			// output.WriteString("\n;;")
			mux.Unlock()
		})

		collUpper.Visit("http://localhost:" + port + "/" + file.Name())

		// fmt.Println(file.Name(), " ", modelNum)
		i := 0
		saveResult := ""
		testCase.testInfo.testResult = append(testCase.testInfo.testResult, "")

		collLower.OnHTML("table.matrix ", func(eTable *colly.HTMLElement) {
			mux.Lock()

			eTable.ForEach("tr", func(_ int, e *colly.HTMLElement) {
				e.ForEach(".cell-full", func(_ int, el *colly.HTMLElement) {
					//output.WriteString("len " + strconv.Itoa(len(testCase.testInfo.testResult[i])) + " ")
					if el.ChildText(".text-success") != "" {
						//testCase.testInfo.testResult[i] = testCase.testInfo.testResult[i] + "1;"
						saveResult = saveResult + "1;"
						// output.WriteString(el.ChildText(".text-success"))
					} else if el.ChildText(".text-light") != "" {
						//testCase.testInfo.testResult[i] = testCase.testInfo.testResult[i] + "x;"
						saveResult = saveResult + "x;"
						// output.WriteString(el.ChildText(".text-light"))
					} else {
						//testCase.testInfo.testResult[i] = testCase.testInfo.testResult[i] + ";"
						saveResult = saveResult + ";"
						//output.WriteString("1 " + testCase.testInfo.testResult[i])
					}
					//output.WriteString(testCase.testInfo.testResult[i])

					// 	el.ForEach(".text-success", func(_ int, _ *colly.HTMLElement) {
					// 		output.WriteString("1")
					// 	})
					// 	el.ForEach(".text-light", func(_ int, _ *colly.HTMLElement) {
					// 		output.WriteString("x")
					// 	})
					// 	output.WriteString(";")
				})

				// output model name
				// body > div > table.matrix.table.table-sm.table-less-padding.table-borderless.table-striped > tbody > tr:nth-child(3) > td.boardmodel.text-ellipsis > span
				e.ForEach(".boardmodel", func(_ int, el *colly.HTMLElement) {
					if el.ChildText("span") != "" {
						testCase.modelInfo.modelName = append(testCase.modelInfo.modelName, el.ChildText("span"))
						//testCase.testInfo.testResult[i] = saveResult
						if i == 0 {
							testCase.testInfo.testResult[0] = saveResult
						} else {
							testCase.testInfo.testResult = append(testCase.testInfo.testResult, saveResult)
						}
						saveResult = ""
						// output.WriteString("\n" + el.ChildText("span") + ";")
						// output.WriteString(testCase.testInfo.testResult[i])

						i++
						//output.WriteString(strconv.Itoa(i))
					}

				})
				// body > div > table.matrix.table.table-sm.table-less-padding.table-borderless.table-striped > tbody > tr:nth-child(11) > td:nth-child(3)
				// str := ""
			})
			mux.Unlock()
		})

		collLower.Visit("http://localhost:" + port + "/" + file.Name())

		outputFileWrite(testCase, output)
		//fmt.Println(output.Name())
		wg.Done()

	}

}

func outputFileWrite(testCase *Info, output *os.File) {
	// output test no.
	for i, _ := range testCase.testInfo.testNo {
		output.WriteString("; " + testCase.testInfo.testNo[i] + " ")
	}
	output.WriteString("\n")

	// output test date
	for i, _ := range testCase.testInfo.testDate {
		output.WriteString("; " + testCase.testInfo.testDate[i] + " ")
	}

	// output excel formula
	output.WriteString("\n;;" + testCase.excelFomula + "\n;;")

	// output test result
	for i, _ := range testCase.modelInfo.modelName {
		output.WriteString("\n" + testCase.modelInfo.modelName[i] + testCase.testInfo.testResult[i])
	}

}

func main() {
	defer os.Exit(0)

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
		testCase := Info{
			testInfo: TestInfo{
				testNum: 0,
			}, modelInfo: ModelInfo{
				modelNum: 0,
			},
		}
		go webClamb(&testCase, outputPath, file, *collUpper, *collLower, port, wg)
	}

	wg.Wait()

}
