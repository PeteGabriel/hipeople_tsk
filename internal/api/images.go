package api

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)




//Upload an image
func (a App) Upload() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}

//GetImage handler function to retrieve an image for a given ID.
func (a App) GetImage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//check info from url
		if match, _ := regexp.MatchString(`/api/image/\d+$`, r.URL.Path); !match {
			//TODO handle not match flow with JSON error
			w.WriteHeader(http.StatusNotFound)
			_, err := w.Write([]byte("image not found"))
			if err != nil {
				//TODO add to logs
				return
			}
			return
		}

		parts := strings.Split(r.URL.Path, "/")
		id, err := strconv.Atoi(parts[3])
		if err != nil {
			fmt.Println("input invalid in url path")
			id = -1
		}

		//TODO send image
		_, err = w.Write([]byte(fmt.Sprintf("%s %d", "Get Image with id ", id)))
		if err != nil {
			//TODO add to logs
			return
		}
	})
}

