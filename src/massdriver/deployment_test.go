package massdriver_test

import (
	"context"
	"testing"
	"xo/src/massdriver"
	testmass "xo/test"
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
			name:   "Test Plan Started",
			id:     "id",
			status: "plan_start",
			want:   `{"metadata":{"timestamp":"2021-01-01 12:00:00.1234","provisioner":"testaform","event_type":"plan_started"},"payload":{"deployment_id":"id"}}`,
		},
		{
			name:   "Test Plan Completed",
			id:     "id",
			status: "plan_complete",
			want:   `{"metadata":{"timestamp":"2021-01-01 12:00:00.1234","provisioner":"testaform","event_type":"plan_completed"},"payload":{"deployment_id":"id"}}`,
		},
		{
			name:   "Test Plan Failed",
			id:     "id",
			status: "plan_fail",
			want:   `{"metadata":{"timestamp":"2021-01-01 12:00:00.1234","provisioner":"testaform","event_type":"plan_failed"},"payload":{"deployment_id":"id"}}`,
		},
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
			testmass.NewMassdriverTestClient("")
			testClient := testmass.NewMassdriverTestClient("")

			err := testClient.MassClient.ReportDeploymentStatus(context.Background(), tc.id, tc.status)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			got := testClient.GetSNSMessages()
			if got[0] != tc.want {
				t.Fatalf("want: %v, got: %v", got, tc.want)
			}
		})
	}
}
