package metadata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnv(t *testing.T) {
	t.Setenv("REGION_ID", "cn-hangzhou")
	t.Setenv("ALIBABA_CLOUD_ACCOUNT_ID", "112233445566")
	t.Setenv("CLUSTER_ID", "c12345678")
	m := &ENVMetadata{}
	expected := map[MetadataKey]string{
		RegionID:  "cn-hangzhou",
		AccountID: "112233445566",
		ClusterID: "c12345678",
	}
	for k, v := range expected {
		value, err := m.Get(k)
		assert.NoError(t, err)
		assert.Equal(t, v, value)
	}

}
