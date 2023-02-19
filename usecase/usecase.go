package usecase

import (
	"context"
	"io"

	"github.com/android-project-46group/core-api/repository"
	"github.com/android-project-46group/core-api/util/logger"
)

type Usecase interface {
	// メンバー関連情報一覧を zip 形式で全取得する。
	DownloadMembersZip(ctx context.Context, writer io.Writer) error
}

type usecase struct {
	database repository.Database
	remote   repository.Remote

	logger logger.Logger
}

func New(database repository.Database, remote repository.Remote, logger logger.Logger) Usecase {
	usecase := &usecase{
		database: database,
		remote:   remote,
		logger:   logger,
	}

	return usecase
}
