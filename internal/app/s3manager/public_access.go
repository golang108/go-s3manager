package s3manager

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// HandleCheckPublicAccess checks if an object is publicly accessible.
func HandleCheckPublicAccess(s3 S3) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bucketName := mux.Vars(r)["bucketName"]
		objectName := mux.Vars(r)["objectName"]

		endpoint := s3.EndpointURL().String()
		if !strings.HasSuffix(endpoint, "/") {
			endpoint += "/"
		}

		// Construct the public URL using path-style access (http://endpoint/bucket/object)
		publicURL := fmt.Sprintf("%s%s/%s", endpoint, bucketName, objectName)

		resp, err := http.Head(publicURL)
		isAccessible := false
		statusCode := 0

		if err != nil {
			isAccessible = false
		} else {
			defer func() { _ = resp.Body.Close() }()
			statusCode = resp.StatusCode
			isAccessible = resp.StatusCode == http.StatusOK
		}

		response := map[string]any{
			"accessible": isAccessible,
			"statusCode": statusCode,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			handleHTTPError(w, fmt.Errorf("error encoding JSON: %w", err))
			return
		}
	}
}
