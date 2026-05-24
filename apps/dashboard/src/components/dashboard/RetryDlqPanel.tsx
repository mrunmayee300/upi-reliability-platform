"use client";

import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { RetryMetrics } from "@/types/dashboard";
import { Archive, RefreshCw, Timer } from "lucide-react";

export function RetryDlqPanel({ retries }: { retries: RetryMetrics }) {
  return (
    <Card>
      <CardHeader>
        <h3 className="text-sm font-semibold text-white">Retry & DLQ</h3>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="flex items-center gap-3">
          <RefreshCw className="h-4 w-4 text-warning" />
          <div>
            <p className="text-xs text-muted">Pending retries</p>
            <p className="text-lg font-semibold text-white">{retries.pending_retries}</p>
          </div>
        </div>
        <div className="flex items-center gap-3">
          <Archive className="h-4 w-4 text-danger" />
          <div>
            <p className="text-xs text-muted">DLQ (24h)</p>
            <p className="text-lg font-semibold text-white">{retries.dlq_count_24h}</p>
          </div>
        </div>
        <div className="flex items-center gap-3">
          <Timer className="h-4 w-4 text-accent" />
          <div>
            <p className="text-xs text-muted">Avg retry attempts</p>
            <p className="text-lg font-semibold text-white">
              {retries.avg_retry_attempts.toFixed(2)}
            </p>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
