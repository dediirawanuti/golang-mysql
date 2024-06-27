package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	"gopkg.in/robfig/cron.v2"
)

type Response struct {
	Status  int64  `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

func IntToString(i int) string {
	s := strconv.Itoa(i)
	return s
}

func Healthz(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Run Healtz")
	var response Response
	response.Status = 200
	response.Message = "OK"
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func StringHash(password string) string {
	h := sha1.New()
	h.Write([]byte(password + os.Getenv("PASSWORD_SALT")))
	sha1_hash := hex.EncodeToString(h.Sum(nil))
	return sha1_hash
}

func Cron() {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	c := cron.New()
	c.AddFunc("* * * * *", func() {
		t1 := time.Now()

		fmt.Println("[Job 1]Every minute job")
		resp, err := client.Get("https://golang-mysql-production.up.railway.app/healthz")

		if err != nil {
			_, _ = http.Get("https://api.telegram.org/bot" + os.Getenv("API_KEY_BOT") + "/sendMessage?chat_id=-1002184332225&text=BE error: " + err.Error())
		} else {
			if resp.StatusCode != 200 {
				_, _ = http.Get("https://api.telegram.org/bot" + os.Getenv("API_KEY_BOT") + "/sendMessage?chat_id=-1002184332225&text=Status code: " + strconv.Itoa(resp.StatusCode))
			}
		}

		t2 := time.Now()
		hs := t2.Sub(t1).Hours()

		hs, mf := math.Modf(hs)
		ms := mf * 60

		ms, sf := math.Modf(ms)
		ss := sf * 60 * 1000
		fmt.Println(ss)

	})
	fmt.Println("Start cron")
	c.Start()
}
