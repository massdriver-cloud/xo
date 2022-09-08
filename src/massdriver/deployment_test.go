package massdriver_test

import (
	"context"
	"testing"
	"xo/src/massdriver"
)

func TestReportDeploymentStatus(t *testing.T) {

	type testData struct {
		name   string
		id     string
		status string
		want   string
	}
	tests := []testData{
		{
			name:   "Test Provision Started",
			id:     "id",
			status: "provision_start",
			want:   `{"metadata":{"timestamp":"2021-01-01 12:00:00.1234","provisioner":"testaform","event_type":"provision_started"},"payload":{"deployment_id":"id"}}`,
		},
		{
			name:   "Test Provision Completed",
			id:     "id",
			status: "provision_complete",
			want:   `{"metadata":{"timestamp":"2021-01-01 12:00:00.1234","provisioner":"testaform","event_type":"provision_completed"},"payload":{"deployment_id":"id"}}`,
		},
		{
			name:   "Test Provision Failed",
			id:     "id",
			status: "provision_fail",
			want:   `{"metadata":{"timestamp":"2021-01-01 12:00:00.1234","provisioner":"testaform","event_type":"provision_failed"},"payload":{"deployment_id":"id"}}`,
		},
		{
			name:   "Test Decommision Started",
			id:     "id",
			status: "decommission_start",
			want:   `{"metadata":{"timestamp":"2021-01-01 12:00:00.1234","provisioner":"testaform","event_type":"decommission_started"},"payload":{"deployment_id":"id"}}`,
		},
		{
			name:   "Test Decommission Completed",
			id:     "id",
			status: "decommission_complete",
			want:   `{"metadata":{"timestamp":"2021-01-01 12:00:00.1234","provisioner":"testaform","event_type":"decommission_completed"},"payload":{"deployment_id":"id"}}`,
		},
		{
			name:   "Test Decommission Failed",
			id:     "id",
			status: "decommission_fail",
			want:   `{"metadata":{"timestamp":"2021-01-01 12:00:00.1234","provisioner":"testaform","event_type":"decommission_failed"},"payload":{"deployment_id":"id"}}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("MASSDRIVER_PROVISIONER", "testaform")
			massdriver.EventTimeString = func() string { return "2021-01-01 12:00:00.1234" }
			testSNSClient := SNSTestClient{}
			testClient := massdriver.MassdriverClient{SNSClient: &testSNSClient, Specification: &massdriver.Specification{}}
			err := testClient.ReportDeploymentStatus(context.Background(), tc.id, tc.status)
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
