package role

import (
	"context"
	"fmt"

	"github.com/fidelfly/gox/errorx"

	"github.com/fidelfly/fxgos/cmd/pkg/db"
	"github.com/fidelfly/fxgos/cmd/service/role/res"
	res2 "github.com/fidelfly/fxgos/cmd/service/user/res"
	"github.com/fidelfly/fxgos/cmd/utilities/mctx"
	"github.com/fidelfly/fxgos/cmd/utilities/pub"
	"github.com/fidelfly/fxgos/cmd/utilities/syserr"
)

func Create(ctx context.Context, resRole *res.Role) (int64, error) {
	if resRole == nil {
		return 0, syserr.ErrInvalidParam
	}
	mctx.FillUserInfo(ctx, resRole)
	if _, err := db.Create(resRole); err != nil {
		return 0, syserr.DatabaseErr(err)
	} else {
		pub.Publish(pub.ResourceEvent{
			Type:   ResourceType,
			Action: pub.ResourceCreate,
			Id:     resRole.Id,
		}, pub.TopicResource)
		return resRole.Id, nil
	}
}

type UpdateInput struct {
	db.UpdateInfo
	Data *res.Role
}

func Update(ctx context.Context, input UpdateInput) error {
	if input.Data == nil {
		return syserr.ErrInvalidParam
	}
	resRole := input.Data

	opts := make([]db.QueryOption, 0)
	if input.Id > 0 {
		resRole.Id = input.Id
		opts = append(opts, db.ID(input.Id))
	} else if resRole.Id > 0 {
		opts = append(opts, db.ID(resRole.Id))
	}
	if len(input.Cols) > 0 {
		opts = append(opts, db.Cols(append(input.Cols, "update_user")...))
	}
	mctx.FillUserInfo(ctx, resRole)
	if rows, err := db.Update(resRole, opts...); err != nil {
		return syserr.DatabaseErr(err)
	} else if rows == 0 {
		return syserr.ErrNotFound
	}
	pub.Publish(pub.ResourceEvent{
		Type:   ResourceType,
		Action: pub.ResourceCreate,
		Id:     resRole.Id,
	}, pub.TopicResource)
	return nil
}

func Read(ctx context.Context, id int64) (*res.Role, error) {
	if id <= 0 {
		return nil, syserr.ErrInvalidParam
	}
	resRole := &res.Role{Id: id}
	if find, err := db.Read(resRole); err != nil {
		return nil, syserr.DatabaseErr(err)
	} else if !find {
		return nil, syserr.ErrNotFound
	}
	return resRole, nil
}

func ReadByCode(ctx context.Context, code string) (*res.Role, error) {
	if len(code) == 0 {
		return nil, syserr.ErrInvalidParam
	}
	resRole := &res.Role{Code: code}
	if find, err := db.Read(resRole); err != nil {
		return nil, syserr.DatabaseErr(err)
	} else if !find {
		return nil, syserr.ErrNotFound
	}
	return resRole, nil
}

func Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return syserr.ErrInvalidParam
	}
	resRole := &res.Role{Id: id}
	if count, err := db.Count(res.Role{}, db.Where("roles like ?", fmt.Sprintf("%%%d%%", id))); err != nil {
		return syserr.DatabaseErr(err)
	} else if count > 0 {
		return errorx.NewError(syserr.CodeOfDatabaseErr, "role_used_by_role")
	}

	if count, err := db.Count(res2.User{}, db.Where("roles like ?", fmt.Sprintf("%%%d%%", id))); err != nil {
		return syserr.DatabaseErr(err)
	} else if count > 0 {
		return errorx.NewError(syserr.CodeOfDatabaseErr, "role_used_by_user")
	}

	if count, err := db.Delete(resRole); err != nil {
		return syserr.DatabaseErr(err)
	} else if count == 0 {
		return syserr.ErrNotFound
	}

	pub.Publish(pub.ResourceEvent{
		Type:   ResourceType,
		Action: pub.ResourceDelete,
		Id:     id,
	}, pub.TopicResource)
	return nil
}

func List(ctx context.Context, input db.ListInfo) ([]*res.Role, int64, error) {
	resRoles := make([]*res.Role, 0)
	opts := make([]db.QueryOption, 0)
	if len(input.Cond) > 0 {
		opts = append(opts, db.Where(input.Cond))
	}
	queOpts := append(db.GetPagingOption(input), opts...)
	if err := db.Find(&resRoles, queOpts...); err != nil {
		return nil, 0, syserr.DatabaseErr(err)
	}

	count := int64(len(resRoles))
	if !(count < input.Results && input.Page == 1) {
		count, _ = db.Count(new(res.Role), opts...)
	}

	return resRoles, count, nil
}

const ResourceType = "resource.role"
