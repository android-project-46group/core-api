package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/android-project-46group/core-api/model"
	"github.com/opentracing/opentracing-go"
)

const ListMembersQuery = `
	SELECT
		m.member_id,
		g.group_name,
		m.name_ja,
		mi.birthday,
		mi.height_cm,
		mi.blood_type,
		mi.generation,
		mi.blog_url,
		mi.img_url,
		m.left_at
	FROM 
		members AS m
	INNER JOIN member_infos AS mi
		ON m.member_id = mi.member_id
		AND mi.locale_id = 1
	INNER JOIN groups AS g
		ON m.group_id = g.group_id;
`

func (d *database) ListMembers(ctx context.Context) ([]*model.Member, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "database.ListMembers")
	defer span.Finish()

	conn, err := d.db.Conn(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to db.Conn: %w", err)
	}

	rows, err := conn.QueryContext(ctx, ListMembersQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to conn.Querytext: %w", err)
	}

	defer func() {
		if err := rows.Close(); err != nil {
			d.logger.Warnf(ctx, "failed to rows.Close: ", err)
		}
	}()

	members := []*model.Member{}

	for rows.Next() {
		var leftAt sql.NullString

		var member model.Member

		err := rows.Scan(
			&member.ID,
			&member.Group,
			&member.Name,
			&member.Birthday,
			&member.Height,
			&member.BloodType,
			&member.Generation,
			&member.BlogURL,
			&member.ImgURL,
			&leftAt,
		)
		if err != nil {
			d.logger.Warnf(ctx, "failed to scan: ", err)

			continue
		}

		if leftAt.Valid {
			member.LeftAt = leftAt.String
		}

		members = append(members, &member)
	}

	return members, nil
}
