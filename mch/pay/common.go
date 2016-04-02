// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package pay

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/chanxuehong/util"
	"github.com/chanxuehong/wechat/mch"
)

// 统一下单.
func UnifiedOrder(pxy *mch.Proxy, req map[string]string) (resp map[string]string, err error) {
	return pxy.PostXML("https://api.mch.weixin.qq.com/pay/unifiedorder", req)
}

// 查询订单.
func OrderQuery(pxy *mch.Proxy, req map[string]string) (resp map[string]string, err error) {
	return pxy.PostXML("https://api.mch.weixin.qq.com/pay/orderquery", req)
}

// 关闭订单.
func CloseOrder(pxy *mch.Proxy, req map[string]string) (resp map[string]string, err error) {
	return pxy.PostXML("https://api.mch.weixin.qq.com/pay/closeorder", req)
}

// 申请退款.
//  NOTE: 请求需要双向证书.
func Refund(pxy *mch.Proxy, req map[string]string) (resp map[string]string, err error) {
	return pxy.PostXML("https://api.mch.weixin.qq.com/secapi/pay/refund", req)
}

// 查询退款.
func RefundQuery(pxy *mch.Proxy, req map[string]string) (resp map[string]string, err error) {
	return pxy.PostXML("https://api.mch.weixin.qq.com/pay/refundquery", req)
}

// 下载对账单到到文件.
func DownloadBill(filepath string, req map[string]string, httpClient *http.Client) (written int64, err error) {
	if req == nil {
		err = errors.New("nil request req")
		return
	}

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

	return downloadBillToWriter(file, req, httpClient)
}

// 下载对账单到 io.Writer.
func DownloadBillToWriter(writer io.Writer, req map[string]string, httpClient *http.Client) (written int64, err error) {
	if writer == nil {
		err = errors.New("nil writer")
		return
	}
	if req == nil {
		err = errors.New("nil request req")
		return
	}
	return downloadBillToWriter(writer, req, httpClient)
}

var (
	// <xml><return_code><![CDATA[FAIL]]></return_code>
	// <return_msg><![CDATA[require POST method]]></return_msg>
	// </xml>
	downloadBillErrorRootNodeStartElement       = []byte("<xml>")
	downloadBillErrorReturnCodeNodeStartElement = []byte("<return_code>")
)

// 下载对账单到 io.Writer.
func downloadBillToWriter(writer io.Writer, req map[string]string, httpClient *http.Client) (written int64, err error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	buf := make([]byte, 32*1024) // 与 io.copyBuffer 里的默认大小一致

	reqBuf := bytes.NewBuffer(buf[:0])
	if err = util.EncodeXMLFromMap(reqBuf, req, "xml"); err != nil {
		return
	}

	httpResp, err := httpClient.Post("https://api.mch.weixin.qq.com/pay/downloadbill", "text/xml; charset=utf-8", reqBuf)
	if err != nil {
		return
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		err = fmt.Errorf("http.Status: %s", httpResp.Status)
		return
	}

	n, err := io.ReadFull(httpResp.Body, buf)
	switch {
	case err == nil:
		// n == len(buf), 可以认为返回的是对账单而不是xml格式的错误信息
		written, err = bytes.NewReader(buf).WriteTo(writer)
		if err != nil {
			return
		}
		var n2 int64
		n2, err = CopyBuffer(writer, httpResp.Body, buf)
		written += n2
		return
	case err == io.ErrUnexpectedEOF:
		readBytes := buf[:n]
		if index := bytes.Index(readBytes, downloadBillErrorRootNodeStartElement); index != -1 {
			if bytes.Contains(readBytes[index+len(downloadBillErrorRootNodeStartElement):], downloadBillErrorReturnCodeNodeStartElement) {
				// 可以认为是错误信息了, 尝试解析xml
				var result mch.Error
				if err = xml.Unmarshal(readBytes, &result); err == nil {
					err = &result
					return
				}
				// err != nil 执行默认的动作, 写入 writer
			}
		}
		return bytes.NewReader(readBytes).WriteTo(writer)
	case err == io.EOF: // 返回空的body
		err = nil
		return
	default: // 其他的错误
		return
	}
}
