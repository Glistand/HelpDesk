"use client";

import { Suspense } from "react";
import { usePathname } from "next/navigation";
import { IconRail } from "./icon-rail";
import { ConversationListPanel } from "./conversation-list-panel";
import { TicketListPanel } from "./ticket-list-panel";
import type { AuthUser, ConversationSummary, TicketSummary } from "@/lib/types";

export function AppShell({
  children,
  conversations,
  tickets,
  user,
  users,
}: {
  children: React.ReactNode;
  conversations: ConversationSummary[];
  tickets: TicketSummary[];
  user: AuthUser | null;
  users: AuthUser[];
}) {
  const pathname = usePathname();
  const showingTickets = pathname.startsWith("/tickets");
  const hideList =
    pathname === "/profile" || pathname === "/team" || pathname === "/design";

  return (
    <div className="flex h-screen overflow-hidden bg-base">
      <IconRail user={user} />
      {!hideList && (
        <Suspense fallback={<ListPanelSkeleton />}>
          {showingTickets ? (
            <TicketListPanel initialTickets={tickets} userId={user?.id || ""} />
          ) : (
            <ConversationListPanel
              initialConversations={conversations}
              users={users}
            />
          )}
        </Suspense>
      )}
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
