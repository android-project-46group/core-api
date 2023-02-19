package usecase

import (
	"context"

	"github.com/android-project-46group/core-api/repository"
)

type Usecase interface {
	// メンバー関連情報一覧を zip 形式で全取得する。
	DownloadMembersZip(ctx context.Context) ([]byte, error)
}

type usecase struct {
	database repository.Database
	remote   repository.Remote
}

func New(database repository.Database, remote repository.Remote) Usecase {
	u := &usecase{
		database: database,
		remote:   remote,
	}

	return u
}
