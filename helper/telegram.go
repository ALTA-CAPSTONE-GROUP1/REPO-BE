package helper

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/app/config"
)

func SendTelegramMessage(phoneNumber string, message string) error {
	token := config.TokenTelegram
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)

	data := url.Values{}
	data.Set("chat_id", phoneNumber) // Nomor tujuan
	data.Set("text", message)        // Pesan yang ingin dikirim

	response, err := http.PostForm(apiURL, data)
	if err != nil {
		return fmt.Errorf("Terjadi kesalahan saat mengirim pesan ke Telegram: %s", err)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		fmt.Println("Pesan berhasil dikirim ke Telegram.")
	} else {
		return fmt.Errorf("Gagal mengirim pesan ke Telegram. Status code: %d", response.StatusCode)
	}

	return nil
}
