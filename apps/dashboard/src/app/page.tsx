"use client";

import { BankHealthGrid } from "@/components/dashboard/BankHealthGrid";
import { Header } from "@/components/dashboard/Header";
import { LiveTransactionFeed } from "@/components/dashboard/LiveTransactionFeed";
import { MetricCards } from "@/components/dashboard/MetricCards";
import { RetryDlqPanel } from "@/components/dashboard/RetryDlqPanel";
import { ThroughputChart } from "@/components/dashboard/ThroughputChart";
import { useLiveDashboard } from "@/hooks/useLiveDashboard";

export default function DashboardPage() {
  const { data, connected, history } = useLiveDashboard();

  return (
    <main className="mx-auto max-w-[1600px] px-4 py-8 sm:px-6 lg:px-8">
      <Header connected={connected} />

      <div className="mt-8 space-y-6">
        <MetricCards summary={data.summary} />

        <div className="grid gap-6 lg:grid-cols-3">
          <ThroughputChart data={history} />
          <RetryDlqPanel retries={data.retries} />
        </div>

        <div className="grid gap-6 lg:grid-cols-2">
          <BankHealthGrid banks={data.banks} />
          <LiveTransactionFeed events={data.recent_events} />
        </div>
      </div>

      <footer className="mt-12 border-t border-border pt-6 text-center text-xs text-muted">
        Phase 4 Dashboard · Analytics WS · Grafana :3001 · Kafka UI :8088 · Jaeger :16686
      </footer>
    </main>
  );
}
