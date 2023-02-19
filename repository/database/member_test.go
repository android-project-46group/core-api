package db_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/android-project-46group/core-api/model"
	db "github.com/android-project-46group/core-api/repository/database"
	"github.com/stretchr/testify/assert"
)

//nolint:funlen
func TestGetPositionsAPI(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name            string
		buildStubs      func(mock sqlmock.Sqlmock)
		expectedMembers []*model.Member
		checkResponse   func(t *testing.T, got []*model.Member, expected []*model.Member, err error)
	}{
		{
			name: "success",
			buildStubs: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectQuery(db.ListMembersQuery).
					WillReturnRows(
						sqlmock.NewRows([]string{
							"member_id", "group_name", "name_ja", "birthday", "height_cm", "blood_type",
							"generation", "blog_url", "img_url", "left_at",
						}).
							AddRow(1, "nogizaka", "john doe", "2004-06-08", 160.0, "O型",
								"4期生", "https://test", "https://test.img", nil).
							AddRow(2, "nogizaka", "john dona", "1999-03-05", 161.0, "O型",
								"2期生", "https://test2", "https://test2.img", "2023-02-19"),
					)
			},
			expectedMembers: []*model.Member{
				{
					ID:         1,
					Group:      "nogizaka",
					Name:       "john doe",
					Birthday:   "2004-06-08",
					Height:     160.0,
					BloodType:  "O型",
					Generation: "4期生",
					BlogURL:    "https://test",
					ImgURL:     "https://test.img",
				},
				{
					ID:         2,
					Group:      "nogizaka",
					Name:       "john dona",
					Birthday:   "1999-03-05",
					Height:     161.0,
					BloodType:  "O型",
					Generation: "2期生",
					BlogURL:    "https://test2",
					ImgURL:     "https://test2.img",
					LeftAt:     "2023-02-19",
				},
			},
			checkResponse: func(t *testing.T, got []*model.Member, expected []*model.Member, err error) {
				t.Helper()

				assert.Nil(t, err)
				assert.Equal(t, expected, got)
			},
		},
		{
			name: "success: empty",
			buildStubs: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectQuery(db.ListMembersQuery).
					WillReturnRows(
						sqlmock.NewRows([]string{
							"member_id", "group_name", "name_ja", "birthday", "height_cm", "blood_type",
							"generation", "blog_url", "img_url", "left_at",
						}),
					)
			},
			expectedMembers: []*model.Member{},
			checkResponse: func(t *testing.T, got []*model.Member, expected []*model.Member, err error) {
				t.Helper()

				assert.Nil(t, err)
				assert.Equal(t, expected, got)
			},
		},
		{
			name: "failed: errConnDone",
			buildStubs: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectQuery(db.ListMembersQuery).
					WillReturnError(sql.ErrConnDone)
			},
			expectedMembers: nil,
			checkResponse: func(t *testing.T, got []*model.Member, expected []*model.Member, err error) {
				t.Helper()

				expectedErrorMessage := fmt.Sprintf("failed to conn.Querytext: %s", sql.ErrConnDone.Error())
				assert.Equal(t, err.Error(), expectedErrorMessage)
				assert.Equal(t, expected, got)
			},
		},
	}

	for _, tc := range testCases {
		//nolint:varnamelen
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			sqlDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				assert.Fail(t, err.Error())
			}

			defer sqlDB.Close()

			tc.buildStubs(mock)
			db := db.NewMockDatabase(sqlDB)

			// Act
			gotMembers, err := db.ListMembers(context.Background())

			// Assert
			tc.checkResponse(t, gotMembers, tc.expectedMembers, err)
		})
	}
}
