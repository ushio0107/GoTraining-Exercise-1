package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/gocolly/colly"
)

type Info struct {
	date      string
	modelName string
	modelFail []int
}

func main() {
	coll := colly.NewCollector()

	data := `./Training_Program/static` // test data path
	files, err := ioutil.ReadDir(data)  // read all the file inside the path above
	if err != nil {
		fmt.Println(err)
	}

	// find "name.csv", if not existed, create one; if existed, clear it.
	// output, err := os.Create("arc.MediaSourceUI" + ".csv")
	// if err != nil {
	// 	fmt.Println("File is failed, err: ", err)
	// }
	// defer output.Close()

	// run a loop find all html file names from ./static
	go func() {
		for _, file := range files {
			if file.IsDir() {
				continue
			} else {
				fmt.Println(file.Name())

				// find "file.name().csv", if not existed, create one; if existed, clear it.
				// the output file will created in path ./output/
				output, err := os.Create("./output/" + file.Name() + ".csv")
				if err != nil {
					fmt.Println("File is failed, err: ", err)
				}
				defer output.Close()
			}
		}
	}()

	// find tag that is class = text-light
	coll.OnHTML(".text-light", func(e *colly.HTMLElement) {
		// fmt.Println(e.Text)

		// output to the file with the same name
		// output.WriteString(e.Text)
	})

	time.Sleep(1 * time.Second)

	coll.Visit("http://localhost:8000/" + "arc.MediaSourceUI.html")

}
