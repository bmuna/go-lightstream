package abr

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type PredictRequest struct {
	BandwidthHistory []float64 `json:"bandwidth_history"`
	BufferSeconds    float64   `json:"buffer_seconds"`
}

type PredictResponse struct {
	RecommendedBitrate int `json:"recommended_bitrate"`
}

func GetRecommendedBitrate(history []float64, buffer float64) (int, error) {
	reqBody := PredictRequest{
		BandwidthHistory: history,
		BufferSeconds:    buffer,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return 0, err
	}

	resp, err := http.Post("http://localhost:5001/predict-bitrate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result PredictResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, err
	}

	return result.RecommendedBitrate, nil
}
