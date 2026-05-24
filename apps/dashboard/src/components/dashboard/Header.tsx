"use client";

import { Radio } from "lucide-react";

export function Header({ connected }: { connected: boolean }) {
  return (
    <header className="flex flex-col gap-4 border-b border-border pb-6 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <p className="text-xs font-medium uppercase tracking-[0.2em] text-accent">
          UPI Platform
        </p>
        <h1 className="mt-1 text-2xl font-bold text-white sm:text-3xl">
          Transaction Intelligence & Failure Recovery
        </h1>
        <p className="mt-1 text-sm text-muted">
          Real-time observability for simulated UPI-scale payment infrastructure
        </p>
      </div>
      <div className="flex items-center gap-2 rounded-full border border-border bg-card px-4 py-2">
        <Radio
          className={`h-4 w-4 ${connected ? "text-success animate-pulse" : "text-danger"}`}
        />
        <span className="text-sm text-white">
          {connected ? "Live" : "Reconnecting…"}
        </span>
      </div>
    </header>
  );
}
