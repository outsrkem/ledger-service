package config

import (
	"fmt"
	"ledger/src/pkg/crypto"
	"ledger/src/slog"
	"os"
)

// encryption 命令行加密
func encryption(plain string) {
	fmt.Println(crypto.Encryption(plain))
	os.Exit(0)
}

// decryCfgCipher 解密配置文件密文
func decryCfgCipher(cipherText *string) {
	klog := slog.FromContext(nil)
	if cipherText != nil && *cipherText == "" {
		return
	}
	plain, err := crypto.Decryption(*cipherText)
	if err != nil {
		klog.Fatal("Decryption password failed. ", cipherText)
		os.Exit(100)
	}
	*cipherText = plain
}
