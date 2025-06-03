package metadata

import (
	"context"
	"strings"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type ProfileMetadata struct {
	profile *v1.ConfigMap
}

var MetadataProfileDataKeys = map[MetadataKey]string{
	ClusterID: "clusterid",
	AccountID: "uid",
}

func NewProfileMetadata(client kubernetes.Interface) (*ProfileMetadata, error) {
	profile, err := client.CoreV1().ConfigMaps("kube-system").Get(context.Background(), "ack-cluster-profile", metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return &ProfileMetadata{profile: profile}, nil
}

func (m *ProfileMetadata) Get(key MetadataKey) (string, error) {
	if key, ok := MetadataProfileDataKeys[key]; ok {
		if m.profile.Data[key] != "" {
			return m.profile.Data[key], nil
		}
	}
	switch key {
	case DataPlaneZoneID:
		vswZone := strings.Split(m.profile.Data["vsw-zone"], ",")
		_, zone, found := strings.Cut(vswZone[0], ":")
		if found {
			return zone, nil
		}
	}
	return "", ErrUnknownMetadataKey
}

type ProfileFetcher struct {
	client kubernetes.Interface
}

func (f *ProfileFetcher) FetchFor(key MetadataKey) (MetadataProvider, error) {
	switch key {
	case DataPlaneZoneID: // supported
	default:
		_, ok := MetadataProfileDataKeys[key]
		if !ok {
			return nil, ErrUnknownMetadataKey
		}
	}
	p, err := NewProfileMetadata(f.client)
	if err != nil {
		return nil, err
	}
	return newImmutableProvider(p, "ClusterProfile"), nil
}
