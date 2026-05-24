"use client";

import { Activity, AlertTriangle, Clock, RefreshCw, Zap } from "lucide-react";
import { Card, CardContent } from "@/components/ui/card";
import { PlatformSummary } from "@/types/dashboard";
import { formatPercent } from "@/lib/utils";

export function MetricCards({ summary }: { summary: PlatformSummary }) {
  const items = [
    {
      label: "TPS",
      value: summary.tps.toFixed(1),
      icon: Zap,
      color: "text-accent",
    },
    {
      label: "Success rate",
      value: formatPercent(summary.success_rate),
      icon: Activity,
      color: "text-success",
    },
    {
      label: "Failure rate",
      value: formatPercent(summary.failure_rate),
      icon: AlertTriangle,
      color: "text-danger",
    },
    {
      label: "Retry rate",
      value: formatPercent(summary.retry_rate),
      icon: RefreshCw,
      color: "text-warning",
    },
    {
      label: "P95 latency",
      value: `${summary.p95_latency_ms} ms`,
      icon: Clock,
      color: "text-accent",
    },
  ];

  return (
    <div className="grid grid-cols-2 gap-4 lg:grid-cols-5">
      {items.map((item) => (
        <Card key={item.label}>
          <CardContent className="pt-5">
            <div className="flex items-start justify-between">
              <div>
                <p className="text-xs font-medium uppercase tracking-wider text-muted">
                  {item.label}
                </p>
                <p className="mt-2 text-2xl font-semibold text-white">{item.value}</p>
              </div>
              <item.icon className={`h-5 w-5 ${item.color}`} />
            </div>
          </CardContent>
        </Card>
      ))}
    </div>
  );
}
