package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func SwaggerRegisterRouter(r *mux.Router) {

	r.Handle("/swagger/swagger.yaml", http.StripPrefix("/swagger/", http.FileServer(http.Dir("./swagger"))))
	r.HandleFunc("/swagger/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		html := `
	<!DOCTYPE html>
	<html>
	<head>
	  <title>Swagger UI</title>
	  <link rel="stylesheet" type="text/css" href="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/4.18.3/swagger-ui.css" />
	</head>
	<body>
	  <div id="swagger-ui"></div>
	  <script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/4.18.3/swagger-ui-bundle.js"></script>
	  <script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/4.18.3/swagger-ui-standalone-preset.js"></script>
	  <script>
	    window.onload = function() {
	      SwaggerUIBundle({
	        url: "/swagger/swagger.yaml",
	        dom_id: '#swagger-ui',
	        presets: [
	          SwaggerUIBundle.presets.apis,
	          SwaggerUIStandalonePreset
	        ],
	        layout: "BaseLayout"
	      });
	    }
	  </script>
	</body>
	</html>
			`
		_, err := w.Write([]byte(html))
		if err != nil {
			log.Error("Could not write response: ", err)
		}
	})

}
