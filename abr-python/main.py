# abr-python/main.py

from fastapi import FastAPI
from pydantic import BaseModel
from predict import predict_bandwidth, select_bitrate

app = FastAPI()

class PredictRequest(BaseModel):
    bandwidth_history: list[float]
    buffer_seconds: float

class PredictResponse(BaseModel):
    recommended_bitrate: int

@app.post("/predict-bitrate", response_model=PredictResponse)
def predict(req: PredictRequest):
    predicted_bw = predict_bandwidth(req.bandwidth_history)
    bitrate = select_bitrate(predicted_bw, req.buffer_seconds)
    return PredictResponse(recommended_bitrate=bitrate)
