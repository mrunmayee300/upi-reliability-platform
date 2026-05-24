"use client";

import { useCallback, useEffect, useRef, useState } from "react";
import { DashboardPayload, emptyDashboard } from "@/types/dashboard";

const WS_URL =
  process.env.NEXT_PUBLIC_WS_URL ?? "ws://localhost:8089/v1/ws/live";
const API_URL =
  process.env.NEXT_PUBLIC_ANALYTICS_URL ?? "http://localhost:8089";

export function useLiveDashboard() {
  const [data, setData] = useState<DashboardPayload>(emptyDashboard());
  const [connected, setConnected] = useState(false);
  const [history, setHistory] = useState<{ t: string; tps: number; p95: number }[]>([]);
  const wsRef = useRef<WebSocket | null>(null);

  const applyPayload = useCallback((payload: DashboardPayload) => {
    setData(payload);
    setHistory((prev) => {
      const point = {
        t: new Date().toLocaleTimeString(),
        tps: payload.summary.tps,
        p95: payload.summary.p95_latency_ms,
      };
      const next = [...prev, point];
      return next.slice(-30);
    });
  }, []);

  useEffect(() => {
    fetch(`${API_URL}/v1/dashboard`)
      .then((r) => r.json())
      .then(applyPayload)
      .catch(() => undefined);

    const connect = () => {
      const ws = new WebSocket(WS_URL);
      wsRef.current = ws;

      ws.onopen = () => setConnected(true);
      ws.onclose = () => {
        setConnected(false);
        setTimeout(connect, 3000);
      };
      ws.onerror = () => ws.close();
      ws.onmessage = (ev) => {
        try {
          applyPayload(JSON.parse(ev.data));
        } catch {
          /* ignore */
        }
      };
    };

    connect();
    return () => wsRef.current?.close();
  }, [applyPayload]);

  return { data, connected, history };
}
