package iam

import (
	"context"

	"github.com/fidelfly/gox/cachex"
	"github.com/fidelfly/gox/logx"

	"github.com/fidelfly/fxgos/cmd/pkg/db"
	"github.com/fidelfly/fxgos/cmd/service/iam/res"
	"github.com/fidelfly/fxgos/cmd/service/user"
	"github.com/fidelfly/fxgos/cmd/utilities/pub"
)

func Initialize() error {
	err := db.Synchronize(
		new(res.Model), new(res.Policy),
	)
	pub.Subscribe(pub.TopicResource, subscriber)
	resDB.CreateJSONIndex("type", "*", "type")
	ScanIam("./assets/iam")
	return err
}

func subscriber(pubData interface{}) error {
	if re, ok := pubData.(pub.ResourceEvent); ok {
		switch re.Type {
		//case role.ResourceType:
		//	return roleChanged(re)
		case user.ResourceType:
			return userChanged(re)
		}
	}
	return nil
}

/*func roleChanged(event pub.ResourceEvent) error {
	switch event.Action {
	case pub.ResourceCreate:
		fallthrough
	case pub.ResourceUpdate:
		break
	case pub.ResourceDelete:
		return DeleteRolePolicy(context.Background(), event.Id)
	}
	return nil
}*/

func userChanged(event pub.ResourceEvent) error {
	switch event.Action {
	case pub.ResourceCreate:
		fallthrough
	case pub.ResourceUpdate:
		logx.Debugf("User[%d] changed, update iam policies now...", event.Id)
		if userData, err := user.Read(context.Background(), event.Id); err == nil {
			return UpdatePolicyByUser(context.Background(), userData.Id, userData.Roles)
		} else {
			return err
		}
	case pub.ResourceDelete:
		return DeleteUserPolicy(context.Background(), event.Id)
	}
	return nil
}

var resDB = cachex.NewBuntCache("./iam_res.db")
