package WebServe

import (
	"fmt"
	"net/http"
)

func WebServer(port string, dataPath string) {
	// fmt.Println("inside Webserver")

	// serve the html file inside ./static
	http.Handle("/", http.FileServer(http.Dir(dataPath)))

	fmt.Println("Listening and Serving port: ", port)

	// port 8000
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println(err)
	}
}

// func rootHandler(w http.ResponseWriter, r *http.Request) {
// 	data := `./static`               // test data path
// 	files, _ := ioutil.ReadDir(data) // read all the file inside the path above
// 	var mu sync.Mutex
// 	i := 1

// 	go func() {
// 		for _, file := range files {
// 			mu.Lock()
// 			if file.IsDir() {
// 				continue
// 			} else {
// 				fmt.Println(i, " ", file.Name())
// 				http.ServeFile(w, r, file.Name())
// 				http.Handle("/", http.FileServer(http.Dir(file.Name()))) //error
// 				err := http.ListenAndServe(":3000", nil)
// 				if err != nil {
// 					fmt.Println(err)
// 				}
// 				i = i + 1
// 			}
// 			mu.Unlock()
// 		}
// 	}()
//}
