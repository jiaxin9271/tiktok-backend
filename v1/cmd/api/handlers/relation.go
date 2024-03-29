package handlers

import (
	"context"
	"errors"
	"github.com/cloudwego/hertz/pkg/app"
	"strconv"
	"tiktok-backend/cmd/api/pack"
	"tiktok-backend/cmd/api/rpc"
	"tiktok-backend/kitex_gen/relation"
	"tiktok-backend/pkg/constants"
	"tiktok-backend/pkg/errno"
)

// RelationAction 登录用户对其他用户进行关注或取消关注。
// 这个很容易出现重复关注的问题
func RelationAction(ctx context.Context, c *app.RequestContext) {
	// 这种方法是postman的body
	//var req pack.RelationActionRequest
	//if err := c.Bind(&req); err != nil {
	//	pack.SendUserRegisterResponse(c, err, -1, "")
	//	return
	//}
	//hlog.Info(req)

	// 这种方法是postman的params
	tokenStr := c.Query("token")
	toUserIdStr := c.Query("to_user_id")
	actionTypeStr := c.Query("action_type") // 1-关注，2-取消关注

	if len(tokenStr) == 0 {
		pack.SendRelationActionResponse(c, errno.ParamErr)
		return
	}

	toUserId, err := strconv.ParseInt(toUserIdStr, 10, 64)
	if err != nil {
		pack.SendRelationActionResponse(c, err)
		return
	}

	actionType, err := strconv.ParseInt(actionTypeStr, 10, 64)
	if err != nil {
		pack.SendRelationActionResponse(c, err)
		return
	}
	if actionType != constants.Follow && actionType != constants.UnFollow {
		pack.SendRelationActionResponse(c, errors.New("actionType error"))
		return
	}

	if err = rpc.RelationAction(context.Background(), &relation.DouyinRelationActionRequest{
		Token:      tokenStr,
		ToUserId:   toUserId,
		ActionType: int32(actionType),
	}); err != nil {
		pack.SendRelationActionResponse(c, err)
		return
	}

	pack.SendRelationActionResponse(c, errno.Success)
}

// RelationFollowList 登录用户关注的所有用户列表。
func RelationFollowList(ctx context.Context, c *app.RequestContext) {
	token := c.Query("token")
	userIdStr := c.Query("user_id")

	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		pack.SendRelationFollowListResponse(c, err, nil)
		return
	}

	userList, err := rpc.RelationFollowList(context.Background(), &relation.DouyinRelationFollowListRequest{
		Token:  token,
		UserId: userId,
	})
	if err != nil {
		pack.SendRelationFollowListResponse(c, err, nil)
		return
	}

	pack.SendRelationFollowListResponse(c, errno.Success, userList)
}

// RelationFollowerList 所有关注登录用户的粉丝列表。
func RelationFollowerList(ctx context.Context, c *app.RequestContext) {
	token := c.Query("token")
	userIdStr := c.Query("user_id")

	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		pack.SendRelationFollowerListResponse(c, err, nil)
		return
	}

	userList, err := rpc.RelationFollowerList(context.Background(), &relation.DouyinRelationFollowerListRequest{
		Token:  token,
		UserId: userId,
	})
	if err != nil {
		pack.SendRelationFollowerListResponse(c, err, nil)
		return
	}

	pack.SendRelationFollowerListResponse(c, errno.Success, userList)
}

// RelationFriendList 所有关注登录用户的粉丝列表。
func RelationFriendList(ctx context.Context, c *app.RequestContext) {
	token := c.Query("token")
	userIdStr := c.Query("user_id")

	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		pack.SendRelationFriendListResponse(c, err, nil)
		return
	}

	userList, err := rpc.RelationFriendList(context.Background(), &relation.DouyinRelationFriendListRequest{
		Token:  token,
		UserId: userId,
	})
	if err != nil {
		pack.SendRelationFriendListResponse(c, err, nil)
		return
	}

	pack.SendRelationFriendListResponse(c, errno.Success, userList)
}
