package durasirawat

import (
	"fmt"
	"net/http"
	"os"
	"bytes"
	"io/ioutil"
	"database/sql"

	"github.com/golang-mysql/scripts/connection"
)

func sendToTelegram(message string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", os.Getenv("API_KEY_BOT"))
	payload := fmt.Sprintf(`{
		"chat_id": "%s",
		"text": "%s"
	}`, os.Getenv("TELEGRAM_CHAT_ID"), message)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request to Telegram: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send message to Telegram: %s", string(body))
	}

	fmt.Println("Pesan berhasil dikirim ke Telegram.")
	return nil
}

func getDurasiRawat(db *sql.DB) (string, error) {
	query := `SELECT kamar_inap.no_rawat, 
                      reg_periksa.no_rkm_medis, reg_periksa.kd_pj,
                      pasien.nm_pasien, 
                      DATEDIFF(CURDATE(), kamar_inap.tgl_masuk) AS durasi
               FROM kamar_inap
               INNER JOIN reg_periksa ON kamar_inap.no_rawat = reg_periksa.no_rawat
               INNER JOIN pasien ON reg_periksa.no_rkm_medis = pasien.no_rkm_medis
               WHERE DATEDIFF(CURDATE(), kamar_inap.tgl_masuk) >= 3
               AND kamar_inap.stts_pulang = '-' AND reg_periksa.kd_pj = 'bpj'`

	rows, err := db.Query(query)
	if err != nil {
		return "", fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	var pesan bytes.Buffer
	pesan.WriteString("Pasien yang dirawat lebih dari 3 hari:\n\n")
	var adaPasien bool

	for rows.Next() {
		adaPasien = true
		var noRawat, noRkmMedis, namaPasien, jenisBayar string
		var durasi int

		if err := rows.Scan(&noRawat, &noRkmMedis, &jenisBayar, &namaPasien, &durasi); err != nil {
			return "", fmt.Errorf("error scanning row: %w", err)
		}

		if jenisBayar == "bpj" {
			jenisBayar = "BPJS"
		}

		pesan.WriteString(fmt.Sprintf("No Rawat       : %s\nNo RM          : %s\nNama           : %s\nLama Dirawat   : %d hari\nJenis Bayar    : %s\n---------------\n", noRawat, noRkmMedis, namaPasien, durasi, jenisBayar))
	}

	if !adaPasien {
		return "Tidak ada pasien yang memenuhi kondisi.", nil
	}

	return pesan.String(), nil
}

func DurasiHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Func Durasi")

	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Connect to the database
	db := connection.Connect()
	if db == nil {
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	message, err := getDurasiRawat(db)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching data: %v", err), http.StatusInternalServerError)
		return
	}

	if err := sendToTelegram(message); err != nil {
		http.Error(w, fmt.Sprintf("Error sending to Telegram: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Notifikasi berhasil dikirim.")
}