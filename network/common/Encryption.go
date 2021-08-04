package common

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"os"
)

func Padding(plainText []byte,blockSize int) []byte{
	//计算要填充的长度
	n:= blockSize-len(plainText)%blockSize
	//对原来的明文填充n个n
	temp:=bytes.Repeat([]byte{byte(n)},n)
	plainText=append(plainText,temp...)
	return plainText
}
//对密文删除填充
func UnPadding(cipherText []byte) []byte{
	//取出密文最后一个字节end
	end:=cipherText[len(cipherText)-1]
	//删除填充
	cipherText=cipherText[:len(cipherText)-int(end)]
	return cipherText
}
//AEC加密（CBC模式）
func AES_CBC_Encrypt(plainText []byte,key []byte) []byte{
	//指定加密算法，返回一个AES算法的Block接口对象
	block,err:=aes.NewCipher(key)
	if err!=nil{
		panic(err)
	}
	//进行填充
	plainText=Padding(plainText,block.BlockSize())
	//指定初始向量vi,长度和block的块尺寸一致
	iv:=[]byte("12345678abcdefgh")
	//指定分组模式，返回一个BlockMode接口对象
	blockMode:=cipher.NewCBCEncrypter(block,iv)
	//加密连续数据库
	cipherText:=make([]byte,len(plainText))
	blockMode.CryptBlocks(cipherText,plainText)
	//返回密文
	return cipherText
}
//AEC解密（CBC模式）
func AES_CBC_Decrypt(cipherText []byte,key []byte) []byte{
	//指定解密算法，返回一个AES算法的Block接口对象
	block,err:=aes.NewCipher(key)
	if err!=nil{
		panic(err)
	}
	//指定初始化向量IV,和加密的一致
	iv:=[]byte("12345678abcdefgh")
	//指定分组模式，返回一个BlockMode接口对象
	blockMode:=cipher.NewCBCDecrypter(block,iv)
	//解密
	plainText:=make([]byte,len(cipherText))
	blockMode.CryptBlocks(plainText,cipherText)
	//删除填充
	plainText=UnPadding(plainText)
	return plainText
}



//RSA生成公私钥
func GenRsaKey(bits int) error {
	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	file, err := os.Create("private.pem")
	if err != nil {
		return err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}
	// 生成公钥文件
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	file, err = os.Create("public.pem")
	if err != nil {
		return err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}
	return nil
}

func RSA_Encrypt(plainText []byte, path string) []byte {
	//打开文件
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	//读取文件的内容
	info, _ := file.Stat()
	buf := make([]byte, info.Size())
	file.Read(buf)
	//pem解码
	block, _ := pem.Decode(buf)
	//x509解码

	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	//类型断言
	publicKey := publicKeyInterface.(*rsa.PublicKey)
	//对明文进行加密
	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, plainText)
	if err != nil {
		panic(err)
	}
	//返回密文
	return cipherText
}

//RSA解密
// cipherText 需要解密的byte数据
// path 私钥文件路径
func RSA_Decrypt(cipherText []byte,path string) []byte{
	//打开文件
	file,err:=os.Open(path)
	if err!=nil{
		panic(err)
	}
	defer file.Close()
	//获取文件内容
	info, _ := file.Stat()
	buf:=make([]byte,info.Size())
	file.Read(buf)
	//pem解码
	block, _ := pem.Decode(buf)
	//X509解码
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err!=nil{
		panic(err)
	}
	//对密文进行解密
	plainText,_:=rsa.DecryptPKCS1v15(rand.Reader,privateKey,cipherText)
	//返回明文
	return plainText
}
//签名
func GetSign(msg []byte,path string)[]byte{
	//取得私钥
	privateKey:=GetRSAPrivateKey(path)
	//计算散列值
	hash := sha256.New()
	hash.Write(msg)
	bytes := hash.Sum(nil)
	//SignPKCS1v15使用RSA PKCS#1 v1.5规定的RSASSA-PKCS1-V1_5-SIGN签名方案计算签名
	sign, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, bytes)
	if err!=nil{
		panic(sign)
	}
	return sign
}
//验证数字签名
func VerifySign(msg []byte,sign []byte,path string)bool{
	//取得公钥
	publicKey:=GetRSAPublicKey(path)
	//计算消息散列值
	hash := sha256.New()
	hash.Write(msg)
	bytes := hash.Sum(nil)
	//验证数字签名
	err:=rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, bytes, sign)
	return err==nil
}
//读取RSA私钥
func GetRSAPrivateKey(path string)*rsa.PrivateKey{
	//读取文件内容
	file, err := os.Open(path)
	if err!=nil{
		panic(err)
	}
	defer file.Close()
	info, _ := file.Stat()
	buf:=make([]byte,info.Size())
	file.Read(buf)
	//pem解码
	block, _ := pem.Decode(buf)
	//X509解码
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	return privateKey
}
//读取RSA公钥
func GetRSAPublicKey(path string)*rsa.PublicKey{
	//读取公钥内容
	file, err := os.Open(path)
	if err!=nil{
		panic(err)
	}
	defer file.Close()
	info, _ := file.Stat()
	buf:=make([]byte,info.Size())
	file.Read(buf)
	//pem解码
	block, _ := pem.Decode(buf)
	//x509解码
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err!=nil{
		panic(err)
	}
	publicKey := publicKeyInterface.(*rsa.PublicKey)
	return publicKey
}
