package s3manager_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudlena/s3manager/internal/app/s3manager"
	"github.com/matryer/is"
)

func TestHandleGetS3Instances(t *testing.T) {
	t.Parallel()

	cases := []struct {
		it            string
		configs       []s3manager.S3InstanceConfig
		expectedIDs   []string
		expectedNames []string
	}{
		{
			it: "returns a single instance",
			configs: []s3manager.S3InstanceConfig{
				{
					Name:            "primary",
					Endpoint:        "localhost:9000",
					AccessKeyID:     "key",
					SecretAccessKey: "secret",
					SignatureType:   "V4",
				},
			},
			expectedIDs:   []string{"1"},
			expectedNames: []string{"primary"},
		},
		{
			it: "returns multiple instances in order",
			configs: []s3manager.S3InstanceConfig{
				{
					Name:            "first",
					Endpoint:        "localhost:9000",
					AccessKeyID:     "key",
					SecretAccessKey: "secret",
					SignatureType:   "V4",
				},
				{
					Name:            "second",
					Endpoint:        "localhost:9001",
					AccessKeyID:     "key2",
					SecretAccessKey: "secret2",
					SignatureType:   "V4",
				},
			},
			expectedIDs:   []string{"1", "2"},
			expectedNames: []string{"first", "second"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.it, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)

			manager, err := s3manager.NewMultiS3Manager(tc.configs)
			is.NoErr(err)

			handler := s3manager.HandleGetS3Instances(manager)
			req, err := http.NewRequest(http.MethodGet, "/api/instances", nil)
			is.NoErr(err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			resp := rr.Result()
			defer func() {
				err = resp.Body.Close()
				is.NoErr(err)
			}()
			body, err := io.ReadAll(resp.Body)
			is.NoErr(err)

			is.Equal(http.StatusOK, resp.StatusCode)

			var result struct {
				Instances []s3manager.S3InstanceInfo `json:"instances"`
			}
			err = json.Unmarshal(body, &result)
			is.NoErr(err)

			is.Equal(len(tc.expectedIDs), len(result.Instances))
			for i, inst := range result.Instances {
				is.Equal(tc.expectedIDs[i], inst.ID)
				is.Equal(tc.expectedNames[i], inst.Name)
			}
		})
	}
}
