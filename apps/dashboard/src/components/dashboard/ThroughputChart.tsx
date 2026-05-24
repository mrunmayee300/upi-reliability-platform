"use client";

import {
  Area,
  AreaChart,
  CartesianGrid,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";
import { Card, CardContent, CardHeader } from "@/components/ui/card";

export function ThroughputChart({
  data,
}: {
  data: { t: string; tps: number; p95: number }[];
}) {
  return (
    <Card className="col-span-2">
      <CardHeader>
        <h3 className="text-sm font-semibold text-white">Throughput & latency</h3>
        <p className="text-xs text-muted">Live window (30 samples)</p>
      </CardHeader>
      <CardContent>
        <div className="h-64">
          <ResponsiveContainer width="100%" height="100%">
            <AreaChart data={data}>
              <defs>
                <linearGradient id="tpsGrad" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor="#3b82f6" stopOpacity={0.4} />
                  <stop offset="95%" stopColor="#3b82f6" stopOpacity={0} />
                </linearGradient>
              </defs>
              <CartesianGrid strokeDasharray="3 3" stroke="#1f2937" />
              <XAxis dataKey="t" tick={{ fill: "#9ca3af", fontSize: 10 }} />
              <YAxis tick={{ fill: "#9ca3af", fontSize: 10 }} />
              <Tooltip
                contentStyle={{
                  background: "#111827",
                  border: "1px solid #1f2937",
                  borderRadius: 8,
                }}
              />
              <Area
                type="monotone"
                dataKey="tps"
                stroke="#3b82f6"
                fill="url(#tpsGrad)"
                name="TPS"
              />
              <Area
                type="monotone"
                dataKey="p95"
                stroke="#10b981"
                fill="none"
                name="P95 ms"
              />
            </AreaChart>
          </ResponsiveContainer>
        </div>
      </CardContent>
    </Card>
  );
}
