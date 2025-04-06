package disbursementservices

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"pg_sandbox/models"
	"pg_sandbox/utils"
)

func CallbackHandler(callbackurl string, payload models.CallbackPayload) {

	jsonData, err := json.Marshal(payload)

	if err != nil {
		utils.Log(slog.LevelError, "Error Mashalling callback payload", "error", err)

		return
	}

	req, err := http.NewRequest("POST", callbackurl, bytes.NewBuffer(jsonData))
	if err != nil {
		utils.Log(slog.LevelError, "Error Mashalling callback payload", "error", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		utils.Log(slog.LevelError, "Error sending callback payload", "error", err)
	}

	defer resp.Body.Close()
	utils.Log(slog.LevelInfo, "Response status", "status", resp.Status)

}
