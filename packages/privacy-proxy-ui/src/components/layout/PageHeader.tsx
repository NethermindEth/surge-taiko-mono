import type { ReactNode } from "react";

interface PageHeaderProps {
  title: string;
  description?: string;
  actions?: ReactNode;
  /** Small inline tag rendered next to the title (e.g. "Read-only"). */
  tag?: ReactNode;
}

export function PageHeader({ title, description, actions, tag }: PageHeaderProps) {
  return (
    <div className="mb-6 flex flex-wrap items-end justify-between gap-3">
      <div>
        <div className="flex items-center gap-3">
          <h1 className="text-2xl font-semibold tracking-tight text-surge-text">
            {title}
          </h1>
          {tag}
        </div>
        {description ? (
          <p className="mt-1 text-sm text-surge-muted">{description}</p>
        ) : null}
      </div>
      {actions ? <div className="flex items-center gap-2">{actions}</div> : null}
    </div>
  );
}
