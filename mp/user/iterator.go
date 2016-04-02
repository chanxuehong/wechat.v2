// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package user

const (
	UserPageSizeLimit = 10000 // 每次拉取的 OPENID 个数最大值为 10000
)

// UserIterator
//
//  iter, err := Client.UserIterator("NextOpenId")
//  if err != nil {
//      // TODO: 增加你的代码
//  }
//
//  for iter.HasNext() {
//      openids, err := iter.NextPage()
//      if err != nil {
//          // TODO: 增加你的代码
//      }
//      // TODO: 增加你的代码
//  }
type UserIterator struct {
	clt *Client // 关联的微信 Client

	lastUserListData  *UserListResult // 最近一次获取的用户数据
	nextPageHasCalled bool            // NextPage() 是否调用过
}

func (iter *UserIterator) TotalCount() int {
	return iter.lastUserListData.TotalCount
}

func (iter *UserIterator) HasNext() bool {
	if !iter.nextPageHasCalled { // 第一次调用需要特殊对待
		return iter.lastUserListData.GotCount > 0 ||
			iter.lastUserListData.NextOpenId != ""
	}

	// 跟文档的描述貌似有点不一样, 返回的 next_openid 并不是列表后一个用户, 而是列表最后一个用户,
	// 并且还要多做最后一次判断才会返回 next_openid=="", 不过这个问题不大, 很多文件读取也是这样, 最后一次返回 0, EOF
	//
	// https://api.weixin.qq.com/cgi-bin/user/get?access_token=k53fhGGYBYCSCEGHj9uveBb9_Y9LigUTtV-L4fJhHuehMCbrtYUsnzVgH9EUejMNNVJLldwLhC81KFEUlInhDO2Zu7KjBXPdDQAgwykW8Go&next_openid=o0seyt4j7FR_tcnOgjK29qoIzZhE
	//
	// 200	OK
	// Connection: keep-alive
	// Date: Mon, 06 Jul 2015 06:16:59 GMT
	// Server: nginx/1.8.0
	// Content-Type: application/json; encoding=utf-8
	// Content-Length: 241
	// {
	//     "total": 6,
	//     "count": 5,
	//     "data": {
	//         "openid": [
	//             "o0seyt0svpBRlwfhv6TGHTgVO6mQ",
	//             "o0seyt1qIaOfmckPrWU-6kCM0oWk",
	//             "o0seytyB45l6wg40Jd8dyAc-Uod0",
	//             "o0seyt8pxzMBI2pttQn5RE9Ce3bk",
	//             "o0seyt0rKDLVlh_VHWPgOORWTe8c"
	//         ]
	//     },
	//     "next_openid": "o0seyt0rKDLVlh_VHWPgOORWTe8c"
	// }
	//
	// https://api.weixin.qq.com/cgi-bin/user/get?access_token=k53fhGGYBYCSCEGHj9uveBb9_Y9LigUTtV-L4fJhHuehMCbrtYUsnzVgH9EUejMNNVJLldwLhC81KFEUlInhDO2Zu7KjBXPdDQAgwykW8Go&next_openid=o0seyt0rKDLVlh_VHWPgOORWTe8c
	//
	// 200	OK
	// Connection: keep-alive
	// Date: Mon, 06 Jul 2015 06:20:30 GMT
	// Server: nginx/1.8.0
	// Content-Type: application/json; encoding=utf-8
	// Content-Length: 38
	// {
	//     "total": 6,
	//     "count": 0,
	//     "next_openid": ""
	// }
	//
	return iter.lastUserListData.NextOpenId != ""
}

func (iter *UserIterator) NextPage() (OpenIdList []string, err error) {
	if !iter.nextPageHasCalled { // 第一次调用需要特殊对待
		iter.nextPageHasCalled = true

		OpenIdList = iter.lastUserListData.Data.OpenIdList
		return
	}

	data, err := iter.clt.UserList(iter.lastUserListData.NextOpenId)
	if err != nil {
		return
	}

	iter.lastUserListData = data

	OpenIdList = data.Data.OpenIdList
	return
}

// 获取用户遍历器, 从 NextOpenId 开始遍历, 如果 NextOpenId == "" 则表示从头遍历.
//  NOTE: 目前微信是从 NextOpenId 下一个用户开始遍历的, 和微信文档描述不一样!!!
func (clt *Client) UserIterator(NextOpenId string) (iter *UserIterator, err error) {
	// 逻辑上相当于第一次调用 UserIterator.NextPage, 因为第一次调用 UserIterator.HasNext 需要数据支撑, 所以提前获取了数据

	data, err := clt.UserList(NextOpenId)
	if err != nil {
		return
	}

	iter = &UserIterator{
		clt:               clt,
		lastUserListData:  data,
		nextPageHasCalled: false,
	}
	return
}
