// refer to  https://www.cloudwego.io/zh/docs/hertz/tutorials/toolkit/toolkit/

namespace go relation_gorm
namespace py relation_gorm
namespace java relation_gorm

enum Code {
     Success         = 0
     ParamInvalid    = 1
     DBErr           = 2
}


struct User {
    1: i64 id
    2: string name
    3: i64 follow_count
    4: i64 follower_count
    5: bool is_follow
}

struct FriendUser {
    1: i64 id
    2: string name
    3: i64 follow_count
    4: i64 follower_count
    5: bool is_follow
    6: string avatar
    7: string message
    8: i64 msgType
}

struct ActionRequest{
    1: i64 user_id (api.body="user_id", api.form="user_id", api.vd="$ > 0")
    2: i64 to_user_id (api.body="to_user_id", api.form="to_user_id", api.vd="$ > 0")
    3: i32 action_type (api.body="action_type", api.form="action_type", api.vd="$ == 1||$ == 2")
}

struct ActionResponse{
   1: Code status_code
   2: string status_msg
}

struct FollowListRequest{
    1: i64 user_id (api.body="user_id", api.form="user_id", api.vd="$ > 0")
}

struct FollowListResponse{
    1: Code status_code
    2: string status_msg
    3: list<User> user_list
}

struct FollowerListRequest{
    1: i64 user_id (api.body="user_id", api.form="user_id", api.vd="$ > 0")
}

struct FollowerListResponse{
    1: Code status_code
    2: string status_msg
    3: list<User> user_list
}

struct FriendListRequest{
    1: i64 user_id (api.body="user_id", api.form="user_id", api.vd="$ > 0")
}

struct FriendListResponse{
    1: Code status_code
    2: string status_msg
    3: list<FriendUser> user_list
}


service UserService {
   ActionResponse Action(1:ActionRequest req)(api.post="/douyin/relation/action/")
   FollowListResponse FollowList(1:FollowListRequest req)(api.post="/douyin/relatioin/follow/list/")
   FollowerListResponse  FollowerList(1: FollowerListRequest req)(api.post="/douyin/relation/follower/list/")
   FriendListResponse FriendList(1:FriendListRequest req)(api.post="/douyin/relation/friend/list/")
}