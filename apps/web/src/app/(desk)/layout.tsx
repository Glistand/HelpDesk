import { AppShell } from "@/components/shell/app-shell";
import { getMe, listConversations, listTickets, listUsers } from "@/lib/api";
import type { AuthUser, ConversationSummary, TicketSummary } from "@/lib/types";

export default async function DeskLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const [conversations, tickets, user] = await Promise.all([
    listConversations().catch((): ConversationSummary[] => []),
    listTickets().catch((): TicketSummary[] => []),
    getMe().catch((): AuthUser | null => null),
  ]);

  let users: AuthUser[] = await listUsers().catch((): AuthUser[] =>
    user ? [user] : [],
  );

  return (
    <AppShell
      conversations={conversations}
      tickets={tickets}
      user={user}
      users={users}
    >
      {children}
    </AppShell>
  );
}
