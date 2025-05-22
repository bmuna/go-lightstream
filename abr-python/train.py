import torch
import torch.nn as nn
import torch.optim as optim
import pandas as pd
from model import GRUBandwidthPredictor

# Load real data from CSV
df = pd.read_csv("bandwidth_data.csv")
data = df["time_kbps"].values.astype(float)

# Normalize
max_val = max(data)
data = data / max_val

# Build input/output sequences
seq_len = 5
X, y = [], []
for i in range(len(data) - seq_len):
    X.append(data[i:i + seq_len])
    y.append(data[i + seq_len])

X = torch.tensor(X, dtype=torch.float32).view(-1, seq_len, 1)
y = torch.tensor(y, dtype=torch.float32).view(-1, 1)

# Define model
model = GRUBandwidthPredictor()
loss_fn = nn.MSELoss()
optimizer = optim.Adam(model.parameters(), lr=0.01)

# Train
for epoch in range(20):
    model.train()
    optimizer.zero_grad()
    output = model(X)
    loss = loss_fn(output, y)
    loss.backward()
    optimizer.step()
    print(f"Epoch {epoch+1}, Loss: {loss.item():.4f}")

# Save
torch.save(model.state_dict(), "gru_model.pth")
print("✅ Model saved to gru_model.pth")
