import { cn } from "@/lib/utils";

const variants: Record<string, string> = {
  healthy: "bg-success/20 text-success border-success/30",
  degraded: "bg-warning/20 text-warning border-warning/30",
  congested: "bg-danger/20 text-danger border-danger/30",
  down: "bg-danger/30 text-danger border-danger/50",
  default: "bg-accent/20 text-accent border-accent/30",
};

export function Badge({
  variant = "default",
  children,
  className,
}: {
  variant?: keyof typeof variants;
  children: React.ReactNode;
  className?: string;
}) {
  return (
    <span
      className={cn(
        "inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-medium",
        variants[variant] ?? variants.default,
        className
      )}
    >
      {children}
    </span>
  );
}
