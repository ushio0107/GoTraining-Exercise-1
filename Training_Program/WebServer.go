package WebServe

import (
	"fmt"
	"net/http"
)

func WebServer(port string, dataPath string) {
	// fmt.Println("inside Webserver")
	fmt.Println("Listening and Serving port: ", port)

	// serve the html file inside ./static
	// srv := httptest.NewServer()
	http.Handle("/", http.FileServer(http.Dir(dataPath)))

	// port 8000
	// httptest.newserver
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println(err)
	}

	// url.parse
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
