package terraform

import "testing"

const sampleShowJSON = `{
  "format_version": "1.0",
  "terraform_version": "1.5.0",
  "values": {
    "root_module": {
      "resources": [
        {
          "address": "aws_db_instance.main", "mode": "managed", "type": "aws_db_instance", "name": "main",
          "values": { "id": "my-db-instance", "name": "production-db", "instance_class": "db.t3.micro" }
        },
        {
          "address": "data.aws_availability_zones.available", "mode": "data", "type": "aws_availability_zones", "name": "available",
          "values": { "id": "us-east-1" }
        }
      ],
      "child_modules": [
        {
          "address": "module.network",
          "resources": [
            { "address": "module.network.aws_security_group.db", "mode": "managed", "type": "aws_security_group", "name": "db",
              "values": { "id": "sg-123456", "name": "db-security-group" } }
          ]
        }
      ]
    }
  }
}`

func TestExtractorResources(t *testing.T) {
	resources, err := Extractor{}.Resources([]byte(sampleShowJSON), map[string]string{"md:instance": "inst-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resources) != 2 {
		t.Fatalf("expected 2 resources (data source excluded), got %d", len(resources))
	}
	if resources[0].Uri != "my-db-instance" || resources[0].Name != "production-db" {
		t.Errorf("unexpected first resource: %+v", resources[0])
	}
	if len(resources[0].Digest["sha256"]) != 64 {
		t.Errorf("expected sha256 config digest, got %q", resources[0].Digest["sha256"])
	}
	if got := resources[0].Annotations.Fields["type"].GetStringValue(); got != "aws:db-instance" {
		t.Errorf("expected type 'aws:db-instance', got %s", got)
	}
	if got := resources[0].Annotations.Fields["md:instance"].GetStringValue(); got != "inst-1" {
		t.Errorf("expected md:instance 'inst-1', got %s", got)
	}
	if resources[1].Uri != "sg-123456" {
		t.Errorf("expected child-module resource 'sg-123456', got %s", resources[1].Uri)
	}
}

func TestExtractorResources_Empty(t *testing.T) {
	resources, err := Extractor{}.Resources([]byte(`{"format_version":"1.0"}`), nil)
	if err != nil || len(resources) != 0 {
		t.Fatalf("expected 0 resources, got %d (err %v)", len(resources), err)
	}
}

func TestExtractorResources_InvalidJSON(t *testing.T) {
	if _, err := (Extractor{}).Resources([]byte(`{bad`), nil); err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestNormalizeType(t *testing.T) {
	cases := map[string]string{
		"aws_db_instance":     "aws:db-instance",
		"google_sql_database": "google:sql-database",
		"nounderscore":        "nounderscore",
	}
	for in, want := range cases {
		if got := normalizeType(in); got != want {
			t.Errorf("normalizeType(%q) = %q, want %q", in, got, want)
		}
	}
}
