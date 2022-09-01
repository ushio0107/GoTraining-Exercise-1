package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/gocolly/colly"
)

var (
	outputPath    string
	dir           string
	outputDirName string

	mux sync.Mutex
)

type Info struct {
	testInfo     TestInfo
	modelInfo    ModelInfo
	excelFormula string
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
	flag.StringVar(&outputPath, "outputPath", ".", "the path of the output directory(the output file will store inside)")
	flag.StringVar(&dir, "dir", ".", "the working directory ")
	flag.StringVar(&outputDirName, "outputDirName", "/output", "the name of the directory which stored the output files")

	flag.Parse()

}

func webClamb(testCase *Info, file os.FileInfo, collUpper colly.Collector, collLower colly.Collector, wg *sync.WaitGroup) {
	// * Host a local server
	// .html file is not allowed to be crawled directly, hosting a local server
	// to make the package colly to implement web crawled successly.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testData_file, err := os.Open(filepath.Join(dir, file.Name()))
		if err != nil {
			fmt.Println(err)
		}
		defer testData_file.Close()
		testData_content, err := ioutil.ReadAll(testData_file)
		if err != nil {
			fmt.Println(err)
		}
		w.Write(testData_content)
	}))

	// To avoid output files aren't exist, os.Create will create file if file doesn't exist,
	// clear file if file already existed, output file name will as the same as input file
	// also, to avoid double file format, .html will be taken out by strings.Replace
	outputFilePath := filepath.Join(outputPath, outputDirName)
	output, err := os.Create(filepath.Join(outputFilePath, strings.Replace(file.Name(), ".html", "", -1)+".csv"))
	fmt.Println("Created output file: ", output.Name())
	if err != nil {
		fmt.Println("File is failed, err: ", err)
	}
	defer output.Close()

	collUpper.OnHTML("table.matrix ", func(eTable *colly.HTMLElement) {
		mux.Lock()

		// Data				: model number
		// Selector Path	: table.matrix > td.boardmodel > span
		// Stored			: testCase.modelInfo.modelNum
		// Usage			: used at excel formula
		eTable.ForEach("td.boardmodel", func(_ int, e *colly.HTMLElement) {
			if e.ChildText("span") != "" {
				testCase.modelInfo.modelNum++
			}
		})

		// Data				: test case No.
		// Selector Path	: table.matrix > div > div:nth-child(1)
		// Stored			: testCase.testInfo.testNo
		eTable.ForEach(".staggered-odd", func(_ int, e *colly.HTMLElement) {
			e.ForEach("div", func(_ int, el *colly.HTMLElement) {
				if el.ChildText("div:nth-child(1)") != "" {
					testCase.testInfo.testNo = append(testCase.testInfo.testNo, el.ChildText("div:nth-child(1)"))

					// Calculate the number of test case, used at excel formula
					testCase.testInfo.testNum++
				}
			})
		})
		// Data				: test case date
		// Selector Path	: table.matrix > tr.title-fg-row > div.staggered-odd > div:nth-child(2)
		eTable.ForEach(".staggered-odd", func(_ int, e *colly.HTMLElement) {
			testCase.testInfo.testDate = append(testCase.testInfo.testDate, e.ChildText("div:nth-child(2)"))
		})

		// Data:Excel Formula
		// stored testCase.excelFormula, var type string
		// * failed
		// =COUNTIF("B9:BtestCaseNum", "x")
		testCase.excelFormula = "\nfailed;123"
		for i := 0; i < testCase.testInfo.testNum; i++ {
			if i < 25 {
				testCase.excelFormula = testCase.excelFormula + "=COUNTIF(" + string(66+i) + "9:" + string(66+i) + strconv.Itoa(testCase.modelInfo.modelNum+8) + ", \"x\");"
			} else {
				testCase.excelFormula = testCase.excelFormula + "=COUNTIF(" + string(65) + string(65+i%25) + "9:" + string(65) + string(65+i%25) + strconv.Itoa(testCase.modelInfo.modelNum+8) + ", \"x\");"
			}
		}

		// * pass
		// =SUM("B9:BtestCaseNum")
		testCase.excelFormula = testCase.excelFormula + "\npass;"
		for i := 0; i < testCase.testInfo.testNum; i++ {
			if i < 25 {
				testCase.excelFormula = testCase.excelFormula + "=SUM(" + string(66+i) + "9:" + string(66+i) + strconv.Itoa(testCase.modelInfo.modelNum+8) + ");"
			} else {
				testCase.excelFormula = testCase.excelFormula + "=SUM(" + string(65) + string(65+i%25) + "9:" + string(65) + string(65+i%25) + strconv.Itoa(testCase.modelInfo.modelNum+8) + ");"
			}
		}

		// * total run
		// =COUNTA("B9:BtestCaseNum")
		testCase.excelFormula = testCase.excelFormula + "\ntotal run;"
		for i := 0; i < testCase.testInfo.testNum; i++ {
			if i < 25 {
				testCase.excelFormula = testCase.excelFormula + "=COUNTA(" + string(66+i) + "9:" + string(66+i) + strconv.Itoa(testCase.modelInfo.modelNum+8) + ");"
			} else {
				testCase.excelFormula = testCase.excelFormula + "=COUNTA(" + string(65) + string(65+i%25) + "9:" + string(65) + string(65+i%25) + strconv.Itoa(testCase.modelInfo.modelNum+8) + ");"
			}
		}

		// * pass rate
		// =IF(B5=0, "N\A", B5/B6)
		testCase.excelFormula = testCase.excelFormula + "\npass rate;"
		for i := 0; i < testCase.testInfo.testNum; i++ {
			if i < 25 {
				testCase.excelFormula = testCase.excelFormula + "=IF(" + string(66+i) + "5=0, \"N/A\"," + string(66+i) + "5/" + string(66+i) + "6);"
			} else {
				testCase.excelFormula = testCase.excelFormula + "=IF(" + string(65) + string(65+i%25) + "5=0, \"N/A\"," + string(65) + string(65+i%25) + "5/" + string(65) + string(65+i%25) + "6);"
			}
		}

		//mux.Unlock()
	})

	// collUpper.Visit("http://localhost:" + port + "/" + file.Name())
	collUpper.Visit((ts.URL))

	collLower.OnHTML("table.matrix ", func(eTable *colly.HTMLElement) {
		//mux.Lock()
		testCaseNum := 0
		saveResult := ""
		testCase.testInfo.testResult = append(testCase.testInfo.testResult, "")

		// Data				: test case result
		// Selector Path	: table.matrix > tr > td.cell-full / table.matrix > tr > td.cell-full > span.text-light / table.matrix > tr > td.cell-full > span.text-success
		// Output			: ";" / "x" / "1"
		eTable.ForEach("tr", func(_ int, e *colly.HTMLElement) {
			e.ForEach(".cell-full", func(_ int, el *colly.HTMLElement) {
				if el.ChildText(".text-success") != "" {
					saveResult = saveResult + "1;"
				} else if el.ChildText(".text-light") != "" {
					saveResult = saveResult + "x;"
				} else {
					saveResult = saveResult + ";"
				}
			})

			// Data				: model name
			// Selector Path	: table.matrix > td.boardmodel> span
			// Stored			: testCase.testInfo.testResult
			e.ForEach(".boardmodel", func(_ int, el *colly.HTMLElement) {
				if el.ChildText("span") != "" {
					testCase.modelInfo.modelName = append(testCase.modelInfo.modelName, el.ChildText("span"))

					// after declared a variable slice string, it is required to store the first data in testResult[0]
					// using append at first will store the first data in testResult[1] instead of testResult[0]
					// so if testCaseNum == 0 can make sure the first data won't store at a wrong slice
					if testCaseNum == 0 {
						testCase.testInfo.testResult[0] = saveResult
					} else {
						testCase.testInfo.testResult = append(testCase.testInfo.testResult, saveResult)
					}
					saveResult = ""
					testCaseNum++
				}
			})

		})
		mux.Unlock()
	})

	collLower.Visit(ts.URL)

	outputFileWrite(testCase, output)
	wg.Done()

}

func outputFileWrite(testCase *Info, output *os.File) {
	// Output all the saved data to the output file ".csv"
	// 1. test case no.
	// 2. test case date
	// 3. the formula calculated failed, pass, total run and pass rate
	// 4. test case result, "x" means fail, "1" means pass, " " means it didn't take test
	// using ";" to seperate all the information
	for i, _ := range testCase.testInfo.testNo {
		output.WriteString("; " + testCase.testInfo.testNo[i] + " ")
	}
	output.WriteString("\n")

	// 2.
	for i, _ := range testCase.testInfo.testDate {
		output.WriteString("; " + testCase.testInfo.testDate[i] + " ")
	}

	// 3.
	output.WriteString("\n;;" + testCase.excelFormula + "\n;;")

	// 4.
	for i, _ := range testCase.modelInfo.modelName {
		output.WriteString("\n" + testCase.modelInfo.modelName[i] + testCase.testInfo.testResult[i])
	}
}

func main() {
	defer os.Exit(0)

	collUpper := colly.NewCollector()
	collLower := colly.NewCollector()
	wg := new(sync.WaitGroup)

	// * Read Files:
	// dir default is ".", so to speak the directory you are in when you run the execution file
	// ioutil.ReadDir to read the file inside the working directory, and declared as "files"
	files, err := ioutil.ReadDir(dir) // read all the file inside the path above
	if err != nil {
		fmt.Println(err)
	}

	// * Output Dir Setting
	// the output files and folder will be created under the test data directory,
	// the files will include in the folder whose name is set by the user, or default "output"
	// the folder will be created if it's not exist.
	err = os.MkdirAll(filepath.Join(outputPath, outputDirName), os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}

	for _, file := range files {
		testCase := Info{
			testInfo: TestInfo{
				testNum: 0,
			}, modelInfo: ModelInfo{
				modelNum: 0,
			},
		}

		if file.IsDir() || !strings.Contains(file.Name(), ".html") {
			// to prevent the input file is not .html file or it's a directory, set a condition if input file type is directory or
			// the file type is not html, pass the file and go to the next one
			continue

		} else {
			// if wg.Add(1) is outside else {}, when there are no file can be processed, the program will struck
			wg.Add(1)
			go webClamb(&testCase, file, *collUpper, *collLower, wg)
		}
	}

	wg.Wait()

}
