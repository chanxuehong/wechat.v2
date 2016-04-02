// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package account

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/chanxuehong/wechat/mp"
)

// 添加客服账号.
//  account:    完整客服账号, 格式为: 账号前缀@公众号微信号, 账号前缀最多10个字符, 必须是英文或者数字字符.
//  nickname:   客服昵称, 最长6个汉字或12个英文字符
//  password:   客服账号登录密码
//  isPwdPlain: 标识 password 是否为明文格式, true 表示是明文密码, false 表示是密文密码.
func AddKfAccount(clt *mp.Client, account, nickname, password string, isPwdPlain bool) (err error) {
	if password == "" {
		return errors.New("empty password")
	}
	if isPwdPlain {
		md5Sum := md5.Sum([]byte(password))
		password = hex.EncodeToString(md5Sum[:])
	}

	request := struct {
		Account  string `json:"kf_account"`
		Nickname string `json:"nickname"`
		Password string `json:"password"`
	}{
		Account:  account,
		Nickname: nickname,
		Password: password,
	}

	var result mp.Error

	incompleteURL := "https://api.weixin.qq.com/customservice/kfaccount/add?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result
		return
	}
	return
}

// 设置客服信息(增量更新, 不更新的可以留空).
//  account:    完整客服账号, 格式为: 账号前缀@公众号微信号, 账号前缀最多10个字符, 必须是英文或者数字字符.
//  nickname:   客服昵称, 最长6个汉字或12个英文字符
//  password:   客服账号登录密码
//  isPwdPlain: 标识 password 是否为明文格式, true 表示是明文密码, false 表示是密文密码.
func SetKfAccount(clt *mp.Client, account, nickname, password string, isPwdPlain bool) (err error) {
	if isPwdPlain && password != "" {
		md5Sum := md5.Sum([]byte(password))
		password = hex.EncodeToString(md5Sum[:])
	}

	request := struct {
		Account  string `json:"kf_account"`
		Nickname string `json:"nickname,omitempty"`
		Password string `json:"password,omitempty"`
	}{
		Account:  account,
		Nickname: nickname,
		Password: password,
	}

	var result mp.Error

	incompleteURL := "https://api.weixin.qq.com/customservice/kfaccount/update?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result
		return
	}
	return
}

// 上传客服头像.
//  开发者可调用本接口来上传图片作为客服人员的头像, 头像图片文件必须是jpg格式, 推荐使用640*640大小的图片以达到最佳效果.
func UploadKfHeadImage(clt *mp.Client, kfAccount, imagePath string) (err error) {
	if kfAccount == "" {
		return errors.New("empty kfAccount")
	}
	file, err := os.Open(imagePath)
	if err != nil {
		return
	}
	defer file.Close()

	return uploadKfHeadImageFromReader(clt, kfAccount, filepath.Base(imagePath), file)
}

// 上传客服头像.
//  开发者可调用本接口来上传图片作为客服人员的头像, 头像图片文件必须是jpg格式, 推荐使用640*640大小的图片以达到最佳效果.
//  注意参数 filename 不是文件路径, 是指定 multipart/form-data 里面文件名称
func UploadKfHeadImageFromReader(clt *mp.Client, kfAccount, filename string, reader io.Reader) (err error) {
	if kfAccount == "" {
		return errors.New("empty kfAccount")
	}
	if filename == "" {
		return errors.New("empty filename")
	}
	if reader == nil {
		return errors.New("nil reader")
	}

	return uploadKfHeadImageFromReader(clt, kfAccount, filename, reader)
}

// 上传客服头像.
//  注意参数 filename 不是文件路径, 是指定 multipart/form-data 里面文件名称
func uploadKfHeadImageFromReader(clt *mp.Client, kfAccount, filename string, reader io.Reader) (err error) {
	var result mp.Error

	// TODO
	//	incompleteURL := "https://api.weixin.qq.com/customservice/kfaccount/uploadheadimg?kf_account=" +
	//		url.QueryEscape(kfAccount) + "&access_token="
	incompleteURL := "https://api.weixin.qq.com/customservice/kfaccount/uploadheadimg?kf_account=" +
		kfAccount + "&access_token="
	fields := []mp.MultipartFormField{{
		ContentType: 0,
		FieldName:   "media",
		FileName:    filename,
		Value:       reader,
	}}
	if err = clt.PostMultipartForm(incompleteURL, fields, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result
		return
	}
	return
}

// 删除客服账号
func DeleteKfAccount(clt *mp.Client, kfAccount string) (err error) {
	var result mp.Error

	// TODO
	//	incompleteURL := "https://api.weixin.qq.com/customservice/kfaccount/del?kf_account=" +
	//		url.QueryEscape(kfAccount) + "&access_token="
	incompleteURL := "https://api.weixin.qq.com/customservice/kfaccount/del?kf_account=" +
		kfAccount + "&access_token="
	if err = clt.GetJSON(incompleteURL, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result
		return
	}
	return
}
