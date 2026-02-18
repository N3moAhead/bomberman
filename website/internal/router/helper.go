package router

import "net/http"

func renderError(err error, w http.ResponseWriter) {
	if err != nil {
		log.Errorln("Render error", err)
		_, err := w.Write([]byte(http.StatusText(500)))
		if err != nil {
			log.Errorln("Error while trying to write status code error", err)
		}
	}
}
