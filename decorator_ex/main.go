package main

import (
	"fmt"
	"web-go-study/decorator_ex/cipher"
	"web-go-study/decorator_ex/lzw"
)

type Component interface {
	Operator(string)
}

var sentData string
var receiveData string

type SendComponent struct {
}

func (self *SendComponent) Operator(data string) {
	sentData = data
}

type ZipComponent struct {
	com Component
}

func (self *ZipComponent) Operator(data string) {
	lzwData, err := lzw.Write([]byte(data))
	if err != nil {
		panic(err)
	}
	self.com.Operator(string(lzwData))
}

type EncryptComponent struct {
	key string
	com Component
}

func (self *EncryptComponent) Operator(data string) {
	encData, err := cipher.Encrypt([]byte(data), self.key)
	if err != nil {
		panic(err)
	}
	self.com.Operator(string(encData))
}

type DecryptComponent struct {
	key string
	com Component
}

func (self *DecryptComponent) Operator(data string) {
	decryptData, err := cipher.Decrypt([]byte(data), self.key)
	if err != nil {
		panic(err)
	}
	self.com.Operator(string(decryptData))
}

type UnzipComponent struct {
	com Component
}

type ReadComponent struct {
}

func (self *ReadComponent) Operator(data string) {
	receiveData = data
}

func (self *UnzipComponent) Operator(data string) {
	unzipData, err := lzw.Read([]byte(data))
	if err != nil {
		panic(err)
	}
	self.com.Operator(string(unzipData))
}

func main() {
	sender := &EncryptComponent{
		key: "abcde",
		com: &ZipComponent{
			com: &SendComponent{}}}

	sender.Operator("Hello World")
	fmt.Println(sentData)

	receiver := &UnzipComponent{
		com: &DecryptComponent{
			key: "abcde",
			com: &ReadComponent{},
		},
	}

	receiver.Operator(sentData)
	fmt.Println(receiveData)
}
