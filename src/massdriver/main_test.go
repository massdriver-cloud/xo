package massdriver

import (
	mocks "xo/src/utils/mocks"
)

func init() {
	Client = &mocks.MockClient{}
}
