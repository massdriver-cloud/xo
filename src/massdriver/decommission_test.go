package massdriver_test

import (
	"context"
	"testing"
	"xo/src/massdriver"
)

func TestReportDecommissionStarted(t *testing.T) {

	type testData struct {
		name string
		id   string
		want string
	}
	tests := []testData{
		{
			name: "Test Decommission Started",
			id:   "id",
			want: `{"metadata":{"timestamp":"2021-01-01 12:00:00.1234","provisioner":"testaform","event_type":"decommission_started"},"payload":{"deployment_id":"id"}}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("MASSDRIVER_PROVISIONER", "testaform")
			massdriver.EventTimeString = func() string { return "2021-01-01 12:00:00.1234" }
			testSNSClient := SNSTestClient{}
			testClient := massdriver.MassdriverClient{SNSClient: &testSNSClient}
			err := testClient.ReportDecommissionStarted(context.Background(), tc.id)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			got := *testSNSClient.Data
			if got != tc.want {
				t.Fatalf("want: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestReportDecommissionCompleted(t *testing.T) {

	type testData struct {
		name string
		id   string
		want string
	}
	tests := []testData{
		{
			name: "Test Decommission Completed",
			id:   "id",
			want: `{"metadata":{"timestamp":"2021-01-01 12:00:00.1234","provisioner":"testaform","event_type":"decommission_completed"},"payload":{"deployment_id":"id"}}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("MASSDRIVER_PROVISIONER", "testaform")
			massdriver.EventTimeString = func() string { return "2021-01-01 12:00:00.1234" }
			testSNSClient := SNSTestClient{}
			testClient := massdriver.MassdriverClient{SNSClient: &testSNSClient}
			err := testClient.ReportDecommissionCompleted(context.Background(), tc.id)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			got := *testSNSClient.Data
			if got != tc.want {
				t.Fatalf("want: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestReportDecommissionFailed(t *testing.T) {

	type testData struct {
		name string
		id   string
		want string
	}
	tests := []testData{
		{
			name: "Test Decommission Failed",
			id:   "id",
			want: `{"metadata":{"timestamp":"2021-01-01 12:00:00.1234","provisioner":"testaform","event_type":"decommission_failed"},"payload":{"deployment_id":"id"}}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("MASSDRIVER_PROVISIONER", "testaform")
			massdriver.EventTimeString = func() string { return "2021-01-01 12:00:00.1234" }
			testSNSClient := SNSTestClient{}
			testClient := massdriver.MassdriverClient{SNSClient: &testSNSClient}
			err := testClient.ReportDecommissionFailed(context.Background(), tc.id)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			got := *testSNSClient.Data
			if got != tc.want {
				t.Fatalf("want: %v, got: %v", tc.want, got)
			}
		})
	}
}
