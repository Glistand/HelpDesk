import { Suspense } from "react";
import { IconRail } from "./icon-rail";
import { TicketListPanel } from "./ticket-list-panel";

export function AppShell({ children }: { children: React.ReactNode }) {
  return (
    <div className="flex h-screen overflow-hidden bg-base">
      <IconRail />
      <Suspense fallback={<ListPanelSkeleton />}>
        <TicketListPanel />
      </Suspense>
      <main className="flex min-w-0 flex-1 flex-col overflow-hidden bg-base">
        {children}
      </main>
    </div>
  );
}

function ListPanelSkeleton() {
  return (
    <aside className="hidden w-[340px] shrink-0 border-r border-border bg-base md:block" />
  );
}
