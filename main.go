package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/kobs30/saiBTC/saibtcgo"
	"github.com/tkanos/gonfig"
)

func main() {
	config_err := gonfig.GetConf("saibtc.config", &btcvalidatorconfig)
	if config_err != nil {
		fmt.Println("Config missed!! ")
		panic(config_err)
	}
	fmt.Println(btcvalidatorconfig)
	srv := &http.Server{
		Addr: btcvalidatorconfig.Host + ":" + btcvalidatorconfig.Port,
	}
	http.HandleFunc("/", api)
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
	//
	//
	//
	// mess := "This is an example of a signed message."
	//
	// btcKey, err := saibtcgo.GenerateKeyPair()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// btcKey.Dump()
	//
	// signature, err := saibtcgo.SignMessage(mess, btcKey.Private)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(signature)
	// valid, err := saibtcgo.VerifySignature(mess, "H9PPoeUtFOli2NwDDgPz2IMRbEhyZ5ngbRrRhsgeOq83CeMmH7tmXCHUmuX6rj0THQPjsMd2K6mBQl6XL8gdAAM=", "1MRBqNJZ5eBQw531YYFYCtp86TMcQQRzYN")
	// // valid, err := saibtcgo.VerifySignature(mess, signature, btcKey.Address)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(valid)
}

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func api(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	r.ParseForm()
	method := strings.Join(r.Form["method"], "")
	switch method {
	case "generateBTC":
		{
			btcKey, err := saibtcgo.GenerateKeyPair()
			if err != nil {
				fmt.Println(err)
				return
			}
			btcKey.Dump()
			w.Write([]byte("{\"Private\":\""))
			w.Write([]byte(btcKey.Private))
			w.Write([]byte("\",\"Public\":\""))
			w.Write([]byte(btcKey.Public))
			w.Write([]byte("\",\"Address\":\""))
			w.Write([]byte(btcKey.Address + "\"}"))
		}
	case "signMessage":
		{
			mess := strings.Join(r.Form["message"], "")
			private := strings.Join(r.Form["p"], "")
			signature, err := saibtcgo.SignMessage(mess, private)
			if err != nil {
				fmt.Println(err)
				return
			}
			w.Write([]byte("{\"message\":\"" + mess + "\",\"signature\":\""))
			w.Write([]byte(signature + "\"}"))
		}
	case "validateSignature":
		{
			mess := strings.Join(r.Form["message"], "")
			address := strings.Join(r.Form["a"], "")
			s := strings.Join(r.Form["signature"], "")
			// https://www.urldecoder.io/golang/
			r := strings.NewReplacer(" ", "+")
			signature := r.Replace(s)
			valid, err := saibtcgo.VerifySignature(mess, signature, address)
			if err != nil {
				// w.Write([]byte("\n"))
				// w.Write([]byte("validation internal error"))
				fmt.Println(err)
				// return
			}
			w.Write([]byte("{\"address\":\"" + address + "\",\"message\":\"" + mess + "\",\"signature\":\""))
			if valid {
				w.Write([]byte("valid"))
			} else {
				w.Write([]byte("not_valid"))
			}
			w.Write([]byte("\"}"))
		}
	case "createAESkey":
		{
			bytes := make([]byte, 32) //generate a random 32 byte key for AES-256
			if _, err := rand.Read(bytes); err != nil {
				panic(err.Error())
			}
			key := hex.EncodeToString(bytes)
			w.Write([]byte(key))
		}
	case "encrypt":
		{
			mess := strings.Join(r.Form["message"], "")
			k := strings.Join(r.Form["k"], "")
			ciphertext := encrypt(mess, k)
			ciphertextBase64 := base64.URLEncoding.EncodeToString([]byte(ciphertext))
			w.Write([]byte(ciphertextBase64))
		}
	case "decrypt":
		{
			cryptoTextBase64 := strings.Join(r.Form["cipher"], "")
			k := strings.Join(r.Form["k"], "")
			encryptedString, _ := base64.URLEncoding.DecodeString(cryptoTextBase64)
			mess := decrypt(string(encryptedString), k)
			w.Write([]byte(mess))
		}
	}
}

type saibtcvalidatorconfig struct {
	Host string
	Port string
}

var btcvalidatorconfig saibtcvalidatorconfig

func encrypt(stringToEncrypt string, keyString string) (encryptedString string) {

	//Since the key is in string, we need to convert decode it to bytes
	key, _ := hex.DecodeString(keyString)
	plaintext := []byte(stringToEncrypt)

	//Create a new Cipher Number from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	//Create a new GCM - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	//https://golang.org/pkg/crypto/cipher/#NewGCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	//Create a nonce. Nonce should be from GCM
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	//Encrypt the data using aesGCM.Seal
	//Since we don't want to save the nonce somewhere else in this case, we add it as a prefix to the encrypted data. The first nonce argument in Seal is the prefix.
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return fmt.Sprintf("%x", ciphertext)
}

func decrypt(encryptedString string, keyString string) (decryptedString string) {

	key, _ := hex.DecodeString(keyString)
	enc, _ := hex.DecodeString(encryptedString)

	//Create a new Cipher Number from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	//Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	//Get the nonce size
	nonceSize := aesGCM.NonceSize()

	//Extract the nonce from the encrypted data
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	//Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}

	return fmt.Sprintf("%s", plaintext)
}
