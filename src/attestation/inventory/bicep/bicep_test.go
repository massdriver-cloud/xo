package bicep

import "testing"

// Mirrors the flattened `az stack ... -o json` CLI blob: top-level resources plus
// other top-level fields (deploymentId, parameters, outputs, …) we ignore.
const sampleStack = `{
  "id": "/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Resources/deploymentStacks/my-stack",
  "name": "my-stack",
  "provisioningState": "succeeded",
  "deploymentId": "/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Resources/deployments/my-stack-1",
  "parameters": { "p1": { "type": "String", "value": "x" } },
  "outputs": { "o1": { "type": "String", "value": "y" } },
  "resources": [
    { "apiVersion": null, "extension": null, "identifiers": null, "type": null, "id": "/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Sql/servers/mysrv/databases/mydb", "resourceGroup": "rg", "status": "managed", "denyStatus": "none" },
    { "apiVersion": null, "extension": null, "identifiers": null, "type": null, "id": "/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/myacct", "resourceGroup": "rg", "status": "managed", "denyStatus": "none" }
  ]
}`

func TestExtractorAssets(t *testing.T) {
	assets, err := Extractor{}.Assets([]byte(sampleStack), map[string]string{"md:instance": "inst-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(assets) != 2 {
		t.Fatalf("expected 2 assets, got %d", len(assets))
	}

	db := assets[0]
	if db.Name != "mydb" {
		t.Errorf("expected name 'mydb', got %s", db.Name)
	}
	if got := db.Annotations.Fields["type"].GetStringValue(); got != "azure:microsoft.sql:servers:databases" {
		t.Errorf("expected nested type, got %s", got)
	}
	if len(db.Digest["sha256"]) != 64 {
		t.Errorf("expected sha256 digest, got %q", db.Digest["sha256"])
	}
	if got := db.Annotations.Fields["md:instance"].GetStringValue(); got != "inst-1" {
		t.Errorf("expected md:instance 'inst-1', got %s", got)
	}

	if got := assets[1].Annotations.Fields["type"].GetStringValue(); got != "azure:microsoft.storage:storageaccounts" {
		t.Errorf("expected storage type, got %s", got)
	}
}

func TestExtractorAssets_PropertiesFallback(t *testing.T) {
	data := `{"properties":{"resources":[{"id":"/subscriptions/s/resourceGroups/rg/providers/Microsoft.Web/sites/app"}]}}`
	assets, err := Extractor{}.Assets([]byte(data), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(assets) != 1 || assets[0].Name != "app" {
		t.Fatalf("expected 1 asset named 'app', got %+v", assets)
	}
}

func TestExtractorAssets_InvalidJSON(t *testing.T) {
	if _, err := (Extractor{}).Assets([]byte(`{bad`), nil); err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
