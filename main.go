package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strings"
	"sync"
)

var (
	urlStore = make(map[string]string)
	mu sync.Mutex
	secretKey = "secretaeskey12345678901234567890"
	lettersRune = []rune(
		"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
)

func encrypt(originalUrl string) string {
	block, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		log.Fatal(err)
	}

	plainText := []byte(originalUrl)
	cipherText := make([]byte, aes.BlockSize+len(plainText))

	iv := cipherText[:aes.BlockSize]

	if _, err := rand.Read(iv); err != nil {
		log.Fatal(err)
	}

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	return hex.EncodeToString(cipherText)
}

func generateShortId() string {
	b := make([]rune, 6)
	for i := range b {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(lettersRune))))
		if err != nil {
			log.Fatal(err)
		}

		b[i] = lettersRune[num.Int64()]
	}

	return string(b)
}


func shortenUrl(w http.ResponseWriter, r *http.Request) {
	originalUrl := r.URL.Query().Get("url")
	if originalUrl == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	if !(strings.HasPrefix(originalUrl, "https://") || strings.HasPrefix(originalUrl, "http://")) {
		http.Error(w, "URL must start with http:// or https://", http.StatusBadRequest)
		return
	}

	encryptedUrl := encrypt(originalUrl)
	shortId := generateShortId()
	
	mu.Lock()
	urlStore[shortId] = encryptedUrl
	mu.Unlock()

	shortUrl := fmt.Sprintf("http://localhost:8080/%s", shortId)
	fmt.Fprintf(w, "Shortened URL: %s\n", shortUrl)
}

func decrypt(encryptedUrl string) string {
	block, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		log.Fatal(err)
	}

	cipherText, err := hex.DecodeString(encryptedUrl)
	if err != nil {
		log.Fatal(err)
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText)
}

func redirectUrl(w http.ResponseWriter, r *http.Request) {
	shortId := r.URL.Path[1:]
	 
	mu.Lock()
	encryptedUrl, ok := urlStore[shortId]
	mu.Unlock()
	if !ok {
		http.Error(w, "Short URL not found", http.StatusNotFound)
		return
	}

	decryptedUrl := decrypt(encryptedUrl)
	http.Redirect(w, r, decryptedUrl, http.StatusFound)
}

func main() {
	http.HandleFunc("/shorten", shortenUrl)
	http.HandleFunc("/", redirectUrl)

	fmt.Println("Server is running on http://localhost:8080/shorten")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	} 
}