import { AppShell } from "@/components/shell/app-shell";
import { listConversations } from "@/lib/api";
import type { ConversationSummary } from "@/lib/types";

export default async function DeskLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  let conversations: ConversationSummary[] = [];
  try {
    conversations = await listConversations();
  } catch {
    conversations = [];
  }
  return <AppShell conversations={conversations}>{children}</AppShell>;
}
