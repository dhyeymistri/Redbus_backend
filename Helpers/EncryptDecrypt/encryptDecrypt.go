package encryptdecrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"
)

// var salt = "iN7el%5"

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func Encrypt(data []byte, passphrase string) []byte {
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext
}

func Decrypt(data []byte, passphrase string) []byte {
	key := []byte(createHash(passphrase))
	// fmt.Println("key: ",key)
	block, err := aes.NewCipher(key)
	if err != nil {
		// fmt.Println("painic1")
		panic(err.Error())
	}
	// fmt.Println("block: ",block)
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		// fmt.Println("painic2")
		panic(err.Error())
	}
	// fmt.Println("gcm: ",gcm)
	nonceSize := gcm.NonceSize()
	// fmt.Println("nonceSize: ",nonceSize)
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	// fmt.Println("nonce: ",nonce)
	// fmt.Println("ciphertext: ",ciphertext)
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		// fmt.Println("painic3")
		panic(err.Error())
	}
	// fmt.Println("plaintext: ",plaintext)
	return plaintext
}

// func main(){
// 	fmt.Println("Pass");
// 	var encryptedValue = Encrypt([]byte("Pass"),"Secret Key");
// 	fmt.Printf("encryptedValue: %x\n", encryptedValue);

// 	var decryptedValue = Decrypt(encryptedValue,"Secret Key");
// 	fmt.Println("decryptedValue", string(decryptedValue));
// }
