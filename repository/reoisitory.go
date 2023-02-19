package repository

import (
	"context"

	"github.com/android-project-46group/core-api/model"
)

// Database 操作を行うための interface。
type Database interface {
	// DB からメンバー詳細情報一覧を取得する。
	ListMembers(ctx context.Context) ([]*model.Member, error)
}

// HTTP 通信を伴う interface。
type Remote interface {
	// 対象 URL の画像を byte 配列で取得する。
	GetImage(ctx context.Context, url string) ([]byte, error)
}
