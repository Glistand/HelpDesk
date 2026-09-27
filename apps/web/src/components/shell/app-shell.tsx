"use client";

import { Suspense } from "react";
import { usePathname } from "next/navigation";
import { IconRail } from "./icon-rail";
import { ConversationListPanel } from "./conversation-list-panel";
import { TicketListPanel } from "./ticket-list-panel";
import type { ConversationSummary, TicketSummary } from "@/lib/types";

export function AppShell({
  children,
  conversations,
  tickets,
}: {
    children: React.ReactNode;
  conversations: ConversationSummary[];
  tickets: TicketSummary[];
}) {
  const pathname = usePathname();
  const showingTickets = pathname.startsWith("/tickets");

  return (
    <div className="flex h-screen overflow-hidden bg-base">
      <IconRail />
      <Suspense fallback={<ListPanelSkeleton />}>
        {showingTickets ? (
          <TicketListPanel initialTickets={tickets} />
        ) : (
          <ConversationListPanel initialConversations={conversations} />
        )}
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
