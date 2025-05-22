# abr-python/predict.py

import torch
import numpy as np
from model import GRUBandwidthPredictor

# Load or define model
model = GRUBandwidthPredictor()
model.load_state_dict(torch.load("gru_model.pth", map_location=torch.device("cpu")))
model.eval()

def predict_bandwidth(history):
    # history: List of float bandwidth samples (e.g., [3200, 3000, 2800])
    input_seq = torch.tensor(history, dtype=torch.float32).view(1, -1, 1)  # [batch, seq_len, input_size]
    with torch.no_grad():
        predicted = model(input_seq)
    return float(predicted.item())

def select_bitrate(predicted_bandwidth, buffer_seconds):
    # Example logic: conservative bitrate selection
    if buffer_seconds < 2:
        safety_margin = 0.5
    elif buffer_seconds < 5:
        safety_margin = 0.7
    else:
        safety_margin = 0.9

    adjusted_bandwidth = predicted_bandwidth * safety_margin

    # Choose from common ABR levels (in kbps)
    bitrates = [144, 360, 720, 1200, 2400, 3600]
    for rate in reversed(bitrates):
        if adjusted_bandwidth >= rate:
            return rate
    return bitrates[0]  # Fallback to lowest
