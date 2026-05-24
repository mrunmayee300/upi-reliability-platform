# Transaction Generator Service

**Port:** 8082 | **Phase 3**

## Responsibilities

- Synthetic UPI transaction generation at configurable TPS
- Realistic payload distribution (banks, amounts, VPAs)
- Configurable failure rate (client-side simulation before bank)
- Load test driver integration (k6 alternative)

## Contract

`shared/contracts/openapi/tx-generator.yaml`
