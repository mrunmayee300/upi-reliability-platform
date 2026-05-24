import os
import random
from datetime import datetime, timedelta, timezone

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel

app = FastAPI(title="AI Prediction Service", version="1.0.0")
app.add_middleware(CORSMiddleware, allow_origins=["*"], allow_methods=["*"])


class CongestionRequest(BaseModel):
    horizon_minutes: int = 15
    bank_codes: list[str] | None = None


@app.get("/health/live")
def liveness():
    return {"status": "ok"}


@app.get("/health/ready")
def readiness():
    return {"status": "ready"}


@app.post("/v1/predict/congestion")
def predict_congestion(req: CongestionRequest):
    banks = req.bank_codes or ["HDFC", "ICICI", "SBI", "AXIS"]
    now = datetime.now(timezone.utc)
    return {
        "forecasts": [
            {
                "bank_code": b,
                "predicted_score": round(random.uniform(0.1, 0.85), 3),
                "confidence": round(random.uniform(0.6, 0.95), 3),
                "peak_at": (now + timedelta(minutes=req.horizon_minutes // 2)).isoformat(),
            }
            for b in banks
        ]
    }


@app.post("/v1/predict/traffic")
def predict_traffic():
    now = datetime.now(timezone.utc)
    return {
        "predicted_tps": round(random.uniform(800, 12000), 1),
        "surge_probability": round(random.uniform(0.1, 0.7), 3),
        "window_start": now.isoformat(),
        "window_end": (now + timedelta(minutes=15)).isoformat(),
    }


@app.get("/v1/recommend/autoscale")
def recommend_autoscale():
    return {
        "service": "tx-ingestion",
        "current_replicas": 3,
        "recommended_replicas": 5,
        "reason": "predicted traffic surge in 15m",
    }


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=int(os.getenv("HTTP_PORT", "8091")))
