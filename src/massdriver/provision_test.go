package massdriver_test

import (
	"testing"
	"xo/src/massdriver"
)

func TestReportProvisionStarted(t *testing.T) {

	type testData struct {
		name   string
		id     string
		client massdriver.MassdriverClient
		want   string
	}
	tests := []testData{
		{
			name:   "Test Provision Started",
			id:     "id",
			client: massdriver.MassdriverClient{Specification: massdriver.Specification{Provisioner: "testaform"}},
			want:   `{"metadata":{"timestamp":"2021-01-01 12:00:00.1234","provisioner":"testaform","event_type":"provision_started"},"payload":{"deployment_id":"id"}}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("MASSDRIVER_PROVISIONER", "testaform")
			massdriver.EventTimeString = func() string { return "2021-01-01 12:00:00.1234" }
			testSNSClient := SNSTestClient{}
			testClient := massdriver.MassdriverClient{SNSClient: &testSNSClient}
			err := testClient.ReportProvisionStarted(tc.id)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			got := *testSNSClient.Data
			if got != tc.want {
				t.Fatalf("want: %v, got: %v", got, tc.want)
			}
		})
	}
}

func TestReportProvisionCompleted(t *testing.T) {

	type testData struct {
		name   string
		id     string
		client massdriver.MassdriverClient
		want   string
	}
	tests := []testData{
		{
			name:   "Test Provision Completed",
			id:     "id",
			client: massdriver.MassdriverClient{Specification: massdriver.Specification{Provisioner: "testaform"}},
			want:   `{"metadata":{"timestamp":"2021-01-01 12:00:00.1234","provisioner":"testaform","event_type":"provision_completed"},"payload":{"deployment_id":"id"}}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("MASSDRIVER_PROVISIONER", "testaform")
			massdriver.EventTimeString = func() string { return "2021-01-01 12:00:00.1234" }
			testSNSClient := SNSTestClient{}
			testClient := massdriver.MassdriverClient{SNSClient: &testSNSClient}
			err := testClient.ReportProvisionCompleted(tc.id)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			got := *testSNSClient.Data
			if got != tc.want {
				t.Fatalf("want: %v, got: %v", got, tc.want)
			}
		})
	}
}

func TestReportProvisionFailed(t *testing.T) {

	type testData struct {
		name   string
		id     string
		client massdriver.MassdriverClient
		want   string
	}
	tests := []testData{
		{
			name:   "Test Provision Failed",
			id:     "id",
			client: massdriver.MassdriverClient{Specification: massdriver.Specification{Provisioner: "testaform"}},
			want:   `{"metadata":{"timestamp":"2021-01-01 12:00:00.1234","provisioner":"testaform","event_type":"provision_failed"},"payload":{"deployment_id":"id"}}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("MASSDRIVER_PROVISIONER", "testaform")
			massdriver.EventTimeString = func() string { return "2021-01-01 12:00:00.1234" }
			testSNSClient := SNSTestClient{}
			testClient := massdriver.MassdriverClient{SNSClient: &testSNSClient}
			err := testClient.ReportProvisionFailed(tc.id)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			got := *testSNSClient.Data
			if got != tc.want {
				t.Fatalf("want: %v, got: %v", got, tc.want)
			}
		})
	}
}
