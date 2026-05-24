"use client";

import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { LiveEvent } from "@/types/dashboard";
import { formatInr, formatTime } from "@/lib/utils";

export function LiveTransactionFeed({ events }: { events: LiveEvent[] }) {
  return (
    <Card className="h-full">
      <CardHeader>
        <h3 className="text-sm font-semibold text-white">Live transaction stream</h3>
        <p className="text-xs text-muted">WebSocket fan-out from Analytics</p>
      </CardHeader>
      <CardContent>
        <div className="max-h-80 space-y-2 overflow-y-auto pr-1">
          {events.length === 0 && (
            <p className="py-8 text-center text-sm text-muted">
              Waiting for events… Start the generator to see traffic.
            </p>
          )}
          {events.map((ev) => (
            <div
              key={`${ev.id}-${ev.timestamp}`}
              className="flex items-center justify-between rounded-lg border border-border bg-background/40 px-3 py-2 text-xs"
            >
              <div className="min-w-0 flex-1">
                <p className="truncate font-mono text-white">{ev.transaction_id}</p>
                <p className="text-muted">
                  {ev.bank_code ?? "—"} · {ev.type}
                </p>
              </div>
              <div className="ml-3 text-right">
                <span
                  className={
                    ev.status === "SUCCESS" ? "text-success" : "text-danger"
                  }
                >
                  {ev.status}
                </span>
                {ev.amount_paise ? (
                  <p className="text-muted">{formatInr(ev.amount_paise)}</p>
                ) : null}
                <p className="text-[10px] text-muted">{formatTime(ev.timestamp)}</p>
              </div>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  );
}
