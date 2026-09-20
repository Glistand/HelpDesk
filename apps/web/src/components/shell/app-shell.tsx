import { Suspense } from "react";
import { IconRail } from "./icon-rail";
import { ConversationListPanel } from "./conversation-list-panel";
import type { ConversationSummary } from "@/lib/types";

export function AppShell({
  children,
  conversations,
}: {
  children: React.ReactNode;
  conversations: ConversationSummary[];
}) {
  return (
    <div className="flex h-screen overflow-hidden bg-base">
      <IconRail />
      <Suspense fallback={<ListPanelSkeleton />}>
        <ConversationListPanel initialConversations={conversations} />
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
