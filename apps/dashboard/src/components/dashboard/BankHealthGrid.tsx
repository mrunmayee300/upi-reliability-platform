"use client";

import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { BankHealth } from "@/types/dashboard";
import { formatPercent } from "@/lib/utils";

function statusVariant(status: string) {
  switch (status) {
    case "HEALTHY":
      return "healthy" as const;
    case "DEGRADED":
      return "degraded" as const;
    case "CONGESTED":
      return "congested" as const;
    default:
      return "down" as const;
  }
}

export function BankHealthGrid({ banks }: { banks: BankHealth[] }) {
  const sorted = [...banks].sort((a, b) => b.error_rate - a.error_rate);

  return (
    <Card>
      <CardHeader>
        <h3 className="text-sm font-semibold text-white">Bank health</h3>
        <p className="text-xs text-muted">Congestion & circuit state</p>
      </CardHeader>
      <CardContent>
        <div className="grid gap-3 sm:grid-cols-2">
          {sorted.map((bank) => (
            <div
              key={bank.bank_code}
              className="rounded-lg border border-border bg-background/50 p-3"
            >
              <div className="flex items-center justify-between">
                <span className="font-mono text-sm font-semibold text-white">
                  {bank.bank_code}
                </span>
                <Badge variant={statusVariant(bank.status)}>{bank.status}</Badge>
              </div>
              <div className="mt-3 grid grid-cols-2 gap-2 text-xs text-muted">
                <div>
                  <span className="block text-[10px] uppercase">Success</span>
                  <span className="text-success">{formatPercent(bank.success_rate)}</span>
                </div>
                <div>
                  <span className="block text-[10px] uppercase">P95</span>
                  <span className="text-white">{bank.p95_latency_ms} ms</span>
                </div>
                <div>
                  <span className="block text-[10px] uppercase">Error</span>
                  <span className="text-danger">{formatPercent(bank.error_rate)}</span>
                </div>
                <div>
                  <span className="block text-[10px] uppercase">Circuit</span>
                  <span className="font-mono text-white">{bank.circuit_state}</span>
                </div>
              </div>
              <div className="mt-2 h-1.5 overflow-hidden rounded-full bg-border">
                <div
                  className="h-full rounded-full bg-gradient-to-r from-success to-warning"
                  style={{ width: `${bank.success_rate * 100}%` }}
                />
              </div>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  );
}
