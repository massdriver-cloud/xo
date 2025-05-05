package artifact_test

type fakeArtifactService struct {
	ShouldError  bool
	CreateCalled bool
	DeleteCalled bool
}
