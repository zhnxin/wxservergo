package wxbizmsgcrypt

import (
	"crypto/sha1"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"time"
)

func GetSHA1(token, timestamp, nonce, encrypt string) string {
	sortlist := []string{token, timestamp, nonce, encrypt}
	sort.Strings(sortlist)
	source := strings.Join(sortlist, "")
	h := sha1.New()
	io.WriteString(h, source)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (msgCrypt *WXBizMsgCrypt) DecryptMsgText(sPostData, sMsgSignature, sTimeStamp, sNonce string) (int, []byte) {
	err, msg := XmlParseExtract(sPostData)
	if err != nil {
		return WXBizMsgCrypt_ParseXml_Error, nil
	}
	return msgCrypt.DecryptMsg(msg, sMsgSignature, sTimeStamp, sNonce)
}

func (msgCrypt *WXBizMsgCrypt) DecryptMsg(msg ReviceMsg, sMsgSignature, sTimeStamp, sNonce string) (int, []byte) {
	signature := GetSHA1(msgCrypt.m_sToken, sTimeStamp, sNonce, msg.Encrypt)
	if sMsgSignature != signature {
		return WXBizMsgCrypt_ValidateSignature_Error, nil
	}
	errCode, plantText := msgCrypt.prpcryptDecrypt(msg.Encrypt)
	if errCode != WXBizMsgCrypt_OK {
		return errCode, plantText
	}
	return WXBizMsgCrypt_OK, plantText
}

func (msgCrypt *WXBizMsgCrypt) VerifyURL(sMsgSignature, sTimeStamp, sNonce, sEchoStr string) (int, []byte) {
	signature := GetSHA1(msgCrypt.m_sToken, sTimeStamp, sNonce, sEchoStr)
	if signature != sMsgSignature {
		return WXBizMsgCrypt_ValidateSignature_Error, nil
	}
	errCode, plantText := msgCrypt.prpcryptDecrypt(sEchoStr)
	if errCode != WXBizMsgCrypt_OK {
		return errCode, plantText
	}
	return WXBizMsgCrypt_OK, plantText

}

func (msgCrypt *WXBizMsgCrypt) EncryptMsg(sReplyMsg, sNonce string) (int, []byte) {
	errCode, ciphertext := msgCrypt.prpcryptEncrypt(sReplyMsg)
	if errCode != WXBizMsgCrypt_OK {
		return errCode, ciphertext
	}
	timeStamp := time.Now().Unix()
	sTimeStamp := strconv.FormatInt(timeStamp, 10)
	encrypt := string(ciphertext[:])
	signature := GetSHA1(msgCrypt.m_sToken, sTimeStamp, sNonce, encrypt)
	ciphertext = []byte(XmlParseGenerate(encrypt, signature, sTimeStamp, sNonce))
	return WXBizMsgCrypt_OK, ciphertext
}
