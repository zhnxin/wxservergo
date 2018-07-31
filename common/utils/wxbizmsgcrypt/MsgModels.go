package wxbizmsgcrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/binary"
	"encoding/xml"
	"fmt"
	"math/rand"
	"time"
)

const (
	WXBizMsgCrypt_OK                      = 0
	WXBizMsgCrypt_ValidateSignature_Error = -40001
	WXBizMsgCrypt_ParseXml_Error          = -40002
	WXBizMsgCrypt_ComputeSignature_Error  = -40003
	WXBizMsgCrypt_IllegalAesKey           = -40004
	WXBizMsgCrypt_ValidateAppid_Error     = -40005
	WXBizMsgCrypt_EncryptAES_Error        = -40006
	WXBizMsgCrypt_DecryptAES_Error        = -40007
	WXBizMsgCrypt_IllegalBuffer           = -40008
	WXBizMsgCrypt_EncodeBase64_Error      = -40009
	WXBizMsgCrypt_DecodeBase64_Error      = -40010
	WXBizMsgCrypt_GenReturnXml_Error      = -40011

	kAesKeySize        = 32
	kAesIVSize         = 16
	kEncodingKeySize   = 43
	kRandEncryptStrLen = 16
	kMsgLen            = 4
	kMaxBase64Size     = 1000000000
)

type ReviceMsg struct {
	ToUserName string `xml:"ToUserName"`
	AgentID    string `xml:"AgentID"`
	Encrypt    string `xml:"Encrypt"`
}

func XmlParseExtract(sPostData string) (error, ReviceMsg) {
	var msg ReviceMsg
	err := xml.Unmarshal([]byte(sPostData), &msg)
	return err, msg
}
func XmlParseGenerate(encrypt, signature, timestamp, nonce string) string {
	return fmt.Sprintf(`<xml>
<Encrypt><![CDATA[%s]]></Encrypt>
<MsgSignature><![CDATA[%s]]></MsgSignature>
<TimeStamp>%s</TimeStamp>
<Nonce><![CDATA[%s]]></Nonce>
</xml>`, encrypt, signature, timestamp, nonce)
}

type WXBizMsgCrypt struct {
	m_sToken  string
	m_sKey    []byte
	m_sCorpID string
}

func GenerateWXBizMsgCrypt(Token, sEncodingAESKey, sCorpID string) (int, WXBizMsgCrypt) {
	aseKey := fmt.Sprintf("%s=", sEncodingAESKey)
	m_sKey, err := base64.StdEncoding.DecodeString(aseKey)
	if err != nil {
		fmt.Println(err)
		return WXBizMsgCrypt_DecodeBase64_Error, WXBizMsgCrypt{}
	}
	if len(m_sKey) != kAesKeySize {
		return WXBizMsgCrypt_IllegalAesKey, WXBizMsgCrypt{}
	}
	return WXBizMsgCrypt_OK, WXBizMsgCrypt{Token, m_sKey, sCorpID}
}

func (msgCrypt *WXBizMsgCrypt) GetCorpID() string {
	return msgCrypt.m_sCorpID
}

func (msgCrypt *WXBizMsgCrypt) getAseKey() []byte {
	return msgCrypt.m_sKey
}
func (msgCrypt *WXBizMsgCrypt) getIV() []byte {
	return msgCrypt.m_sKey[:kAesIVSize]
}

func PKCS7Decode(decrypted []byte) []byte {
	length := len(decrypted)
	pad := uint8(decrypted[length-1])
	if pad < 1 || pad > 32 {
		pad = 0
	}
	return decrypted[0 : length-int(pad)]
}

func (msgCrypt *WXBizMsgCrypt) prpcryptDecrypt(encryptText string) (int, []byte) {
	ciphertext, err := base64.StdEncoding.DecodeString(encryptText)
	if err != nil {
		fmt.Println(err)
		return WXBizMsgCrypt_DecodeBase64_Error, nil
	}
	block, err := aes.NewCipher(msgCrypt.getAseKey())
	if err != nil {
		return WXBizMsgCrypt_IllegalAesKey, nil
	}

	blockModel := cipher.NewCBCDecrypter(block, msgCrypt.getIV())
	plantText := make([]byte, len(ciphertext))
	blockModel.CryptBlocks(plantText, ciphertext)
	//remove padding
	plantText = PKCS7Decode(plantText)
	networkOrder := plantText[kRandEncryptStrLen : kRandEncryptStrLen+kMsgLen]
	xmlLength := recoverNetworkBytesOrder(networkOrder)
	xmlContent := plantText[kRandEncryptStrLen+kMsgLen : kRandEncryptStrLen+kMsgLen+xmlLength]
	fromCorpID := string(plantText[kRandEncryptStrLen+kMsgLen+xmlLength:])

	if fromCorpID == msgCrypt.m_sCorpID {
		return WXBizMsgCrypt_OK, xmlContent
	} else {
		return WXBizMsgCrypt_ValidateAppid_Error, xmlContent
	}
}

func recoverNetworkBytesOrder(orderBytes []byte) uint32 {
	b_buf := bytes.NewBuffer(orderBytes)
	var x uint32
	binary.Read(b_buf, binary.BigEndian, &x)
	return x
}
func getNetworkBytesOrder(sourceNumber int) []byte {
	orderBytes := make([]byte, 4)
	orderBytes[3] = (byte)(sourceNumber & 0xFF)
	orderBytes[2] = (byte)(sourceNumber >> 8 & 0xFF)
	orderBytes[1] = (byte)(sourceNumber >> 16 & 0xFF)
	orderBytes[0] = (byte)(sourceNumber >> 24 & 0xFF)
	return orderBytes
}
func getRandomStr() string {
	kinds, result := [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}, make([]byte, kRandEncryptStrLen)
	var ikind int
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < kRandEncryptStrLen; i++ {
		// random ikind
		ikind = rand.Intn(3)
		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return string(result)
}
func bytesCombine(pBytes ...[]byte) []byte {
	return bytes.Join(pBytes, []byte(""))
}
func (msgCrypt *WXBizMsgCrypt) generateEncrypt(plaintext string) []byte {
	textBytes := []byte(plaintext)
	ciphertext := bytesCombine([]byte(getRandomStr()),
		getNetworkBytesOrder(len(textBytes)),
		textBytes,
		[]byte(msgCrypt.m_sCorpID))
	return ciphertext
}

func (msgCrypt *WXBizMsgCrypt) prpcryptEncrypt(plaintext string) (int, []byte) {
	ciphertext := msgCrypt.generateEncrypt(plaintext)
	block, err := aes.NewCipher(msgCrypt.getAseKey())
	if err != nil {
		fmt.Println(err)
		return WXBizMsgCrypt_EncryptAES_Error, nil
	}
	//PKCS7Encode
	padding := block.BlockSize() - len(ciphertext)%block.BlockSize()
	if padding == 0 {
		padding = block.BlockSize()
	}
	padtext := bytes.Repeat([]byte{byte(padding)}, padding) //生成填充的文本
	ciphertext = append(ciphertext, padtext...)
	//CBC encrypt
	blockModel := cipher.NewCBCEncrypter(block, msgCrypt.getIV())
	encrypted := make([]byte, len(ciphertext))
	blockModel.CryptBlocks(encrypted, ciphertext)
	encryptedText := base64.StdEncoding.EncodeToString(encrypted)
	return WXBizMsgCrypt_OK, []byte(encryptedText)

}
