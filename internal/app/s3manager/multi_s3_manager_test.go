package s3manager_test

import (
	"testing"

	"github.com/cloudlena/s3manager/internal/app/s3manager"
	"github.com/matryer/is"
)

func TestNewMultiS3Manager(t *testing.T) {
	t.Parallel()

	cases := []struct {
		it          string
		configs     []s3manager.S3InstanceConfig
		expectError bool
	}{
		{
			it: "creates manager with a single V4 instance",
			configs: []s3manager.S3InstanceConfig{
				{
					Name:            "test",
					Endpoint:        "localhost:9000",
					AccessKeyID:     "key",
					SecretAccessKey: "secret",
					SignatureType:   "V4",
				},
			},
		},
		{
			it: "creates manager with V2 signature type",
			configs: []s3manager.S3InstanceConfig{
				{
					Name:            "test",
					Endpoint:        "localhost:9000",
					AccessKeyID:     "key",
					SecretAccessKey: "secret",
					SignatureType:   "V2",
				},
			},
		},
		{
			it: "creates manager with Anonymous signature type",
			configs: []s3manager.S3InstanceConfig{
				{
					Name:          "test",
					Endpoint:      "localhost:9000",
					SignatureType: "Anonymous",
				},
			},
		},
		{
			it: "creates manager with V4Streaming signature type",
			configs: []s3manager.S3InstanceConfig{
				{
					Name:            "test",
					Endpoint:        "localhost:9000",
					AccessKeyID:     "key",
					SecretAccessKey: "secret",
					SignatureType:   "V4Streaming",
				},
			},
		},
		{
			it: "creates manager with IAM credentials",
			configs: []s3manager.S3InstanceConfig{
				{
					Name:     "test",
					Endpoint: "localhost:9000",
					UseIam:   true,
				},
			},
		},
		{
			it: "creates manager with region",
			configs: []s3manager.S3InstanceConfig{
				{
					Name:            "test",
					Endpoint:        "s3.amazonaws.com",
					AccessKeyID:     "key",
					SecretAccessKey: "secret",
					SignatureType:   "V4",
					Region:          "us-east-1",
				},
			},
		},
		{
			it: "creates manager with multiple instances",
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
		},
		{
			it:          "returns error for empty configs",
			configs:     []s3manager.S3InstanceConfig{},
			expectError: true,
		},
		{
			it: "returns error for invalid signature type",
			configs: []s3manager.S3InstanceConfig{
				{
					Name:          "test",
					Endpoint:      "localhost:9000",
					SignatureType: "INVALID",
				},
			},
			expectError: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.it, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)

			manager, err := s3manager.NewMultiS3Manager(tc.configs)
			if tc.expectError {
				is.True(err != nil)
				is.True(manager == nil)
			} else {
				is.NoErr(err)
				is.True(manager != nil)
			}
		})
	}
}

func TestMultiS3ManagerGetInstance(t *testing.T) {
	t.Parallel()

	configs := []s3manager.S3InstanceConfig{
		{
			Name:            "primary",
			Endpoint:        "localhost:9000",
			AccessKeyID:     "key",
			SecretAccessKey: "secret",
			SignatureType:   "V4",
		},
		{
			Name:            "secondary",
			Endpoint:        "localhost:9001",
			AccessKeyID:     "key2",
			SecretAccessKey: "secret2",
			SignatureType:   "V4",
		},
	}

	manager, err := s3manager.NewMultiS3Manager(configs)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("find by ID", func(t *testing.T) {
		t.Parallel()
		is := is.New(t)

		instance, err := manager.GetInstance("1")
		is.NoErr(err)
		is.Equal("primary", instance.Name)
		is.Equal("1", instance.ID)
	})

	t.Run("find by name", func(t *testing.T) {
		t.Parallel()
		is := is.New(t)

		instance, err := manager.GetInstance("secondary")
		is.NoErr(err)
		is.Equal("secondary", instance.Name)
		is.Equal("2", instance.ID)
	})

	t.Run("returns error for unknown identifier", func(t *testing.T) {
		t.Parallel()
		is := is.New(t)

		_, err := manager.GetInstance("does-not-exist")
		is.True(err != nil)
	})
}

func TestMultiS3ManagerGetAllInstances(t *testing.T) {
	t.Parallel()
	is := is.New(t)

	configs := []s3manager.S3InstanceConfig{
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
	}

	manager, err := s3manager.NewMultiS3Manager(configs)
	is.NoErr(err)

	instances := manager.GetAllInstances()
	is.Equal(2, len(instances))
	is.Equal("first", instances[0].Name)
	is.Equal("second", instances[1].Name)
}

func TestMultiS3ManagerGetCurrentClient(t *testing.T) {
	t.Parallel()
	is := is.New(t)

	configs := []s3manager.S3InstanceConfig{
		{
			Name:            "primary",
			Endpoint:        "localhost:9000",
			AccessKeyID:     "key",
			SecretAccessKey: "secret",
			SignatureType:   "V4",
		},
	}

	manager, err := s3manager.NewMultiS3Manager(configs)
	is.NoErr(err)

	client := manager.GetCurrentClient()
	is.True(client != nil)
}

func TestMultiS3ManagerGetCurrentInstance(t *testing.T) {
	t.Parallel()
	is := is.New(t)

	configs := []s3manager.S3InstanceConfig{
		{
			Name:            "primary",
			Endpoint:        "localhost:9000",
			AccessKeyID:     "key",
			SecretAccessKey: "secret",
			SignatureType:   "V4",
		},
	}

	manager, err := s3manager.NewMultiS3Manager(configs)
	is.NoErr(err)

	instance := manager.GetCurrentInstance()
	is.True(instance != nil)
	is.Equal("primary", instance.Name)
}

func TestMultiS3ManagerSetCurrentInstance(t *testing.T) {
	t.Parallel()
	is := is.New(t)

	configs := []s3manager.S3InstanceConfig{
		{
			Name:            "primary",
			Endpoint:        "localhost:9000",
			AccessKeyID:     "key",
			SecretAccessKey: "secret",
			SignatureType:   "V4",
		},
	}

	manager, err := s3manager.NewMultiS3Manager(configs)
	is.NoErr(err)

	err = manager.SetCurrentInstance("1")
	is.True(err != nil)
}
