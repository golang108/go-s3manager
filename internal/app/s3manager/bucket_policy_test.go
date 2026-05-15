package s3manager_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cloudlena/s3manager/internal/app/s3manager"
	"github.com/cloudlena/s3manager/internal/app/s3manager/mocks"
	"github.com/gorilla/mux"
	"github.com/matryer/is"
)

func TestHandleGetBucketPolicy(t *testing.T) {
	t.Parallel()

	cases := []struct {
		it                   string
		getBucketPolicyFunc  func(context.Context, string) (string, error)
		expectedStatusCode   int
		expectedBodyContains string
	}{
		{
			it: "returns the bucket policy",
			getBucketPolicyFunc: func(context.Context, string) (string, error) {
				return `{"Version":"2012-10-17"}`, nil
			},
			expectedStatusCode:   http.StatusOK,
			expectedBodyContains: `{"Version":"2012-10-17"}`,
		},
		{
			it: "returns empty body when policy is empty",
			getBucketPolicyFunc: func(context.Context, string) (string, error) {
				return "", nil
			},
			expectedStatusCode:   http.StatusOK,
			expectedBodyContains: "",
		},
		{
			it: "returns error if there is an S3 error",
			getBucketPolicyFunc: func(context.Context, string) (string, error) {
				return "", errS3
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedBodyContains: "mocked s3 error",
		},
	}

	for _, tc := range cases {
		t.Run(tc.it, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)

			s3 := &mocks.S3Mock{
				GetBucketPolicyFunc: tc.getBucketPolicyFunc,
			}

			r := mux.NewRouter()
			r.Handle("/api/buckets/{bucketName}/policy", s3manager.HandleGetBucketPolicy(s3)).Methods(http.MethodGet)

			ts := httptest.NewServer(r)
			defer ts.Close()

			resp, err := http.Get(ts.URL + "/api/buckets/my-bucket/policy")
			is.NoErr(err)
			defer func() {
				err = resp.Body.Close()
				is.NoErr(err)
			}()
			body, err := io.ReadAll(resp.Body)
			is.NoErr(err)

			is.Equal(tc.expectedStatusCode, resp.StatusCode)
			is.True(strings.Contains(string(body), tc.expectedBodyContains))
		})
	}
}

func TestHandlePutBucketPolicy(t *testing.T) {
	t.Parallel()

	cases := []struct {
		it                   string
		setBucketPolicyFunc  func(context.Context, string, string) error
		body                 string
		expectedStatusCode   int
		expectedBodyContains string
	}{
		{
			it: "sets the bucket policy",
			setBucketPolicyFunc: func(context.Context, string, string) error {
				return nil
			},
			body:               `{"Version":"2012-10-17"}`,
			expectedStatusCode: http.StatusNoContent,
		},
		{
			it: "returns error if there is an S3 error",
			setBucketPolicyFunc: func(context.Context, string, string) error {
				return errS3
			},
			body:                 `{"Version":"2012-10-17"}`,
			expectedStatusCode:   http.StatusInternalServerError,
			expectedBodyContains: "mocked s3 error",
		},
	}

	for _, tc := range cases {
		t.Run(tc.it, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)

			s3 := &mocks.S3Mock{
				SetBucketPolicyFunc: tc.setBucketPolicyFunc,
			}

			r := mux.NewRouter()
			r.Handle("/api/buckets/{bucketName}/policy", s3manager.HandlePutBucketPolicy(s3)).Methods(http.MethodPut)

			ts := httptest.NewServer(r)
			defer ts.Close()

			req, err := http.NewRequest(http.MethodPut, ts.URL+"/api/buckets/my-bucket/policy", bytes.NewBufferString(tc.body))
			is.NoErr(err)

			resp, err := http.DefaultClient.Do(req)
			is.NoErr(err)
			defer func() {
				err = resp.Body.Close()
				is.NoErr(err)
			}()
			body, err := io.ReadAll(resp.Body)
			is.NoErr(err)

			is.Equal(tc.expectedStatusCode, resp.StatusCode)
			is.True(strings.Contains(string(body), tc.expectedBodyContains))
		})
	}
}
