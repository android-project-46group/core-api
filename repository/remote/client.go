package remote

import (
	"net/http"

	"github.com/android-project-46group/core-api/repository"
)

type remote struct {
	client *http.Client
}

func New() repository.Remote {
	//nolint:exhaustivestruct,exhaustruct
	client := &http.Client{
		Transport: newTransport(),
	}
	rm := &remote{
		client: client,
	}

	return rm
}
