/*
 Copyright 2015 Crunchy Data Solutions, Inc.
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package sec

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"github.com/crunchydata/crunchy-postgresql-manager/logit"
)

//stuff used by the crypto routines
//The key argument should be the AES key, either 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256.
var key = "opensesame123456" // 16 bytes!
var ciphertext = []byte("abcdef1234567890")

func EncryptPassword(inputPassword string) (string, error) {
	var err error
	var encryptedRaw []byte

	stuff := []byte(inputPassword)

	encryptedRaw, err = basicencrypt(stuff)
	if err != nil {
		logit.Error.Println(err.Error())
		return "", err
	}

	var encryptedHexString = ""
	//convert the raw encrypted bytes to a hex string for storing in PG
	for i := 0; i < len(encryptedRaw); i++ {
		//fmt.Printf("%x", encryptedRaw[i])
		encryptedHexString = encryptedHexString + fmt.Sprintf("%x", encryptedRaw[i])
	}

	return encryptedHexString, nil

}

func DecryptPassword(encodedHexPassword string) (string, error) {
	var encryptedPassword []byte
	var err error
	var unencryptedPassword string

	encryptedPassword, err = hex.DecodeString(encodedHexPassword)
	if err != nil {
		logit.Error.Println(err.Error())
		return "", err
	}

	unencryptedPassword, err = basicdecrypt(encryptedPassword)
	if err != nil {
		logit.Error.Println(err.Error())
		return "", err
	}
	return unencryptedPassword, nil

}

func basicencrypt(input []byte) ([]byte, error) {
	var output []byte
	block, err := aes.NewCipher([]byte(key))

	if err != nil {
		return output, err
	}

	// 16 bytes for AES-128, 24 bytes for AES-192, 32 bytes for AES-256
	iv := ciphertext[:aes.BlockSize] // const BlockSize = 16

	// encrypt

	encrypter := cipher.NewCFBEncrypter(block, iv)

	encrypted := make([]byte, len(input))
	encrypter.XORKeyStream(encrypted, input)
	//var stroutput = string(encrypted[:])
	//logit.Info.Println("encrypted value " + stroutput)

	return encrypted, nil
}

//expects an encrypted input string
func basicdecrypt(input []byte) (string, error) {

	block, err := aes.NewCipher([]byte(key))

	if err != nil {
		return "", err
	}

	// 16 bytes for AES-128, 24 bytes for AES-192, 32 bytes for AES-256
	iv := ciphertext[:aes.BlockSize] // const BlockSize = 16

	decrypter := cipher.NewCFBDecrypter(block, iv) // simple!

	decrypted := make([]byte, len(input))
	decrypter.XORKeyStream(decrypted, input)

	var stroutput = string(decrypted[:])
	//logit.Info.Println("decrypted value " + stroutput)
	return stroutput, nil
}
