package usecase_test

import (
	"context"
	"io"

	"github.com/android-project-46group/core-api/model"
)

type mockDatabase struct {
	ListMembersFunc func(ctx context.Context) ([]*model.Member, error)
}

func (m *mockDatabase) ListMembers(ctx context.Context) ([]*model.Member, error) {
	return m.ListMembersFunc(ctx)
}

type mockRemote struct {
	GetImageFunc func(ctx context.Context, url string) (io.ReadCloser, error)
}

func (m *mockRemote) GetImage(ctx context.Context, url string) (io.ReadCloser, error) {
	return m.GetImageFunc(ctx, url)
}
