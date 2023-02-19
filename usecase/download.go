package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/android-project-46group/core-api/model"
	"github.com/google/uuid"
)

const (
	workingDir = "tmp"
	imgDir     = "imgs"
)

func (u *usecase) DownloadMembersZip(ctx context.Context) ([]byte, error) {
	randomPath, err := u.createUniqueDir(workingDir)
	if err != nil {
		return nil, fmt.Errorf("failed to createUniqueDir: %w", err)
	}

	members, err := u.database.ListMembers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ListMembers: %w", err)
	}

	jsonPath := filepath.Join(randomPath, "members_info.json")
	if err := u.writeMembersJSON(members, jsonPath); err != nil {
		return nil, fmt.Errorf("faild to writeMembersJson: %w", err)
	}

	if err := u.initializeImgDir(filepath.Join(randomPath, imgDir)); err != nil {
		return nil, fmt.Errorf("faild to initializeImgDir: %w", err)
	}

	for _, member := range members {
		fileName := path.Base(member.ImgURL)
		fullPath := filepath.Join(randomPath, imgDir, fileName)

		err := u.downloadImage(ctx, member.ImgURL, fullPath)
		if err != nil {
			u.logger.Warnf(ctx, "failed to downloadImage: %v", err)
		}
	}

	return nil, nil
}

func (u *usecase) writeMembersJSON(members []*model.Member, fullPath string) error {
	file, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("failed to os.Create: %w", err)
	}
	defer file.Close()

	bytes, err := json.MarshalIndent(members, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to json.MarshalIndent: %w", err)
	}

	_, err = file.Write(bytes)
	if err != nil {
		return fmt.Errorf("failed to f.Write: %w", err)
	}

	return nil
}

func (u *usecase) downloadImage(ctx context.Context, url, fullPath string) error {
	reader, closer, err := u.remote.GetImage(ctx, url)
	if err != nil {
		return fmt.Errorf("failed to GetImage: %w", err)
	}

	defer closer()

	file, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("failed to os.Create: %w", err)
	}

	defer file.Close()

	_, err = io.Copy(file, reader)
	if err != nil {
		return fmt.Errorf("failed to io.Copy: %w", err)
	}

	return nil
}

// 引数で受け取った basePath 配下に、unique なフォルダを作成する。
//
// フォルダ名は uuid で構成され、return の string では
// basePath も含んだパス全体が返される。
func (u *usecase) createUniqueDir(basePath string) (string, error) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("failed to NewRandom: %w", err)
	}

	uuidStr := uuid.String()
	fullPath := filepath.Join(basePath, uuidStr)

	err = os.MkdirAll(fullPath, 0o777)
	if err != nil {
		return "", fmt.Errorf("failed to MkdirAll: %w", err)
	}

	return fullPath, nil
}

func (u *usecase) initializeImgDir(imgDir string) error {
	if err := os.MkdirAll(imgDir, 0o777); err != nil {
		return fmt.Errorf("failed to MkdirAll: %w", err)
	}

	return nil
}
