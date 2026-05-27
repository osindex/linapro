// This file verifies active-release manifest/config resource projection for
// dynamic-plugin host services.

package runtime

import (
	"testing"

	"lina-core/internal/service/plugin/internal/catalog"
)

func TestBuildArtifactResourceViews(t *testing.T) {
	manifest := &catalog.Manifest{
		RuntimeArtifact: &catalog.ArtifactSpec{
			ManifestResources: []*catalog.ArtifactManifestResource{
				{
					Path:    "manifest/config/config.example.yaml",
					Content: []byte("template: true\n"),
				},
				{
					Path:    "manifest/config/config.yaml",
					Content: []byte("runtime: true\n"),
				},
				{
					Path:    "manifest/metadata.yaml",
					Content: []byte("name: demo\n"),
				},
				{
					Path:    "manifest/resources/policy.yaml",
					Content: []byte("enabled: true\n"),
				},
			},
		},
	}

	defaultConfig := buildArtifactDefaultConfig(manifest)
	if string(defaultConfig) != "runtime: true\n" {
		t.Fatalf("expected runtime config only, got %q", string(defaultConfig))
	}

	resources := buildArtifactManifestResources(manifest)
	if len(resources) != 2 {
		t.Fatalf("expected two declaration resources, got %#v", resources)
	}
	if string(resources["metadata.yaml"]) != "name: demo\n" {
		t.Fatalf("expected metadata resource, got %#v", resources)
	}
	if string(resources["resources/policy.yaml"]) != "enabled: true\n" {
		t.Fatalf("expected policy resource, got %#v", resources)
	}
	if _, ok := resources["config/config.yaml"]; ok {
		t.Fatalf("expected config resources to be excluded from manifest view, got %#v", resources)
	}
}
