# AI Prediction Service

**Port:** 8091 | **Language:** Python (FastAPI) | **Phase 3**

## Responsibilities

- Ingest time-series features from `analytics-events`, `bank-health`
- XGBoost / Prophet congestion forecasting
- Traffic surge prediction
- Autoscaling recommendations for HPA
- Publish forecast `congestion-events`

## Stack

- FastAPI, scikit-learn, XGBoost or Prophet
- Model training script: `services/ai-prediction/training/` (Phase 3)

## Contract

`shared/contracts/openapi/ai-prediction.yaml`
