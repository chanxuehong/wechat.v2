// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package media

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"

	"github.com/chanxuehong/wechat/mp"
)

// 下载多媒体到文件.
//  请注意, 视频文件不支持下载
func (clt *Client) DownloadMedia(mediaId, filepath string) (written int64, err error) {
	file, err := os.Create(filepath)
	if err != nil {
		return
	}
	defer func() {
		file.Close()
		if err != nil {
			os.Remove(filepath)
		}
	}()

	return clt.downloadMediaToWriter(mediaId, file)
}

// 下载多媒体到 io.Writer.
//  请注意, 视频文件不支持下载
func (clt *Client) DownloadMediaToWriter(mediaId string, writer io.Writer) (written int64, err error) {
	if writer == nil {
		err = errors.New("nil writer")
		return
	}
	return clt.downloadMediaToWriter(mediaId, writer)
}

// 下载多媒体到 io.Writer.
func (clt *Client) downloadMediaToWriter(mediaId string, writer io.Writer) (written int64, err error) {
	token, err := clt.Token()
	if err != nil {
		return
	}

	hasRetried := false
RETRY:
	finalURL := "https://api.weixin.qq.com/cgi-bin/media/get?media_id=" + url.QueryEscape(mediaId) +
		"&access_token=" + url.QueryEscape(token)

	httpResp, err := clt.HttpClient.Get(finalURL)
	if err != nil {
		return
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		err = fmt.Errorf("http.Status: %s", httpResp.Status)
		return
	}

	ContentType, _, _ := mime.ParseMediaType(httpResp.Header.Get("Content-Type"))
	if ContentType != "text/plain" && ContentType != "application/json" { // 返回的是媒体流
		return io.Copy(writer, httpResp.Body)
	}

	// 返回的是错误信息
	var result mp.Error
	if err = json.NewDecoder(httpResp.Body).Decode(&result); err != nil {
		return
	}

	switch result.ErrCode {
	case mp.ErrCodeOK:
		return // 基本不会出现
	case mp.ErrCodeInvalidCredential, mp.ErrCodeAccessTokenExpired: // 失效(过期)重试一次
		mp.LogInfoln("[WECHAT_RETRY] err_code:", result.ErrCode, ", err_msg:", result.ErrMsg)
		mp.LogInfoln("[WECHAT_RETRY] current token:", token)

		if !hasRetried {
			hasRetried = true

			if token, err = clt.TokenRefresh(); err != nil {
				return
			}
			mp.LogInfoln("[WECHAT_RETRY] new token:", token)

			result = mp.Error{}
			goto RETRY
		}
		mp.LogInfoln("[WECHAT_RETRY] fallthrough, current token:", token)
		fallthrough
	default:
		err = &result
		return
	}
}

// 创建图文消息素材.
func (clt *Client) CreateNews(articles []Article) (info *MediaInfo, err error) {
	if len(articles) <= 0 {
		err = errors.New("图文素材是空的")
		return
	}
	if len(articles) > NewsArticleCountLimit {
		err = fmt.Errorf("图文素材的文章个数不能超过 %d, 现在为 %d", NewsArticleCountLimit, len(articles))
		return
	}

	var request = struct {
		Articles []Article `json:"articles,omitempty"`
	}{
		Articles: articles,
	}

	var result struct {
		mp.Error
		MediaInfo
	}

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/media/uploadnews?access_token="
	if err = ((*mp.Client)(clt)).PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}
	info = &result.MediaInfo
	return
}

// 创建视频素材.
//  mediaId:     通过上传视频文件得到
//  title:       标题, 可以为空
//  description: 描述, 可以为空
func (clt *Client) CreateVideo(mediaId, title, description string) (info *MediaInfo, err error) {
	if mediaId == "" {
		err = errors.New("empty mediaId")
		return
	}
	var request = struct {
		MediaId     string `json:"media_id"`
		Title       string `json:"title,omitempty"`
		Description string `json:"description,omitempty"`
	}{
		MediaId:     mediaId,
		Title:       title,
		Description: description,
	}

	var result struct {
		mp.Error
		MediaInfo
	}

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/media/uploadvideo?access_token="
	if err = ((*mp.Client)(clt)).PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}
	info = &result.MediaInfo
	return
}
