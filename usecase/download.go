package usecase

import (
	"archive/zip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/android-project-46group/core-api/model"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

const (
	workingDir   = "tmp"
	imgDir       = "imgs"
	jsonFileName = "members_info.json"
)

func (u *usecase) DownloadMembersZip(ctx context.Context, writer io.Writer) error {
	randomPath, err := u.createUniqueDir(workingDir)
	if err != nil {
		return fmt.Errorf("failed to createUniqueDir: %w", err)
	}

	// 処理終了時に一時ファイルを削除する。
	defer func() {
		if err := os.RemoveAll(randomPath); err != nil {
			u.logger.Warnf(ctx, "failed to remove dir: %v", err)
		}
	}()

	members, err := u.database.ListMembers(ctx)
	if err != nil {
		return fmt.Errorf("failed to ListMembers: %w", err)
	}

	jsonPath := filepath.Join(randomPath, jsonFileName)
	if err := u.writeMembersJSON(members, jsonPath); err != nil {
		return fmt.Errorf("faild to writeMembersJson: %w", err)
	}

	if err := u.initializeImgDir(filepath.Join(randomPath, imgDir)); err != nil {
		return fmt.Errorf("faild to initializeImgDir: %w", err)
	}

	//nolint:varnamelen
	eg := errgroup.Group{}
	imgParallel := make(chan bool, len(members))

	for _, member := range members {
		member := member
		imgParallel <- true

		eg.Go(func() error {
			fileName := path.Base(member.ImgURL)
			fullPath := filepath.Join(randomPath, imgDir, fileName)

			err := u.downloadImage(ctx, member.ImgURL, fullPath)
			if err != nil {
				u.logger.Warnf(ctx, "failed to downloadImage: %v", err)
			}

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return fmt.Errorf("failed to image parallel fetch: %w", err)
	}

	close(imgParallel)

	err = u.createZip(randomPath, writer)
	if err != nil {
		u.logger.Warnf(ctx, "failed to createZip: %v", err)
	}

	return nil
}

func (u *usecase) writeMembersJSON(members []*model.Member, fullPath string) (err error) {
	file, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("failed to os.Create: %w", err)
	}

	defer func() {
		closeErr := file.Close()
		if closeErr != nil {
			err = errors.Join(err, fmt.Errorf("failed to close json file: %w", closeErr))
		}
	}()

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

func (u *usecase) downloadImage(ctx context.Context, url, fullPath string) (err error) {
	readCloser, err := u.remote.GetImage(ctx, url)
	if err != nil {
		return fmt.Errorf("failed to GetImage: %w", err)
	}

	defer func() {
		closeErr := readCloser.Close()
		if closeErr != nil {
			err = errors.Join(err, fmt.Errorf("failed to close readCloser: %w", closeErr))
		}
	}()

	file, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("failed to os.Create: %w", err)
	}

	defer func() {
		closeErr := file.Close()
		if closeErr != nil {
			err = errors.Join(err, fmt.Errorf("failed to close img file: %w", closeErr))
		}
	}()

	_, err = io.Copy(file, readCloser)
	if err != nil {
		return fmt.Errorf("failed to io.Copy: %w", err)
	}

	return nil
}

// 引数で受け取った basePath 配下に、unique なフォルダを作成する。
//
// フォルダ名は uuid であり、return の string では basePath も含んだパス全体が返される。
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

// targetDir 配下のファイルを再帰的に zip 化し、渡された writer に書き込む。
func (u *usecase) createZip(targetDir string, readWriter io.Writer) (err error) {
	zipWriter := zip.NewWriter(readWriter)
	defer func() {
		closeErr := zipWriter.Close()
		if closeErr != nil {
			err = errors.Join(err, fmt.Errorf("failed to close img file: %w", closeErr))
		}
	}()

	// ディレクトリを再帰的に探索する。
	err = filepath.Walk(targetDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to Walk: %w", err)
		}

		// 圧縮するファイル名を生成する。
		relPath, err := filepath.Rel(targetDir, path)
		if err != nil {
			return fmt.Errorf("failed to Rel: %w", err)
		}
		zipPath := filepath.ToSlash(filepath.Join(filepath.Dir(relPath), filepath.Base(path)))

		// ディレクトリは飛ばす。
		if info.IsDir() {
			return nil
		}

		// 圧縮する。
		fileToZip, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("failed to Open: %w", err)
		}

		defer func() {
			closeErr := fileToZip.Close()
			if closeErr != nil {
				err = errors.Join(err, closeErr)
			}
		}()

		fileToZipStat, err := fileToZip.Stat()
		if err != nil {
			return fmt.Errorf("failed to Stat: %w", err)
		}

		header, err := zip.FileInfoHeader(fileToZipStat)
		if err != nil {
			return fmt.Errorf("failed to FileInfoHeader: %w", err)
		}

		header.Name = zipPath
		header.Method = zip.Deflate

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return fmt.Errorf("failed to CreateHeader: %w", err)
		}

		if _, err = io.Copy(writer, fileToZip); err != nil {
			return fmt.Errorf("failed to io.Copy: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to %w", err)
	}

	return nil
}
