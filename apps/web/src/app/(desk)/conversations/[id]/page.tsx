import { notFound } from "next/navigation";
import { ConversationThread } from "@/components/conversation-thread";
import { ApiError, getConversation, listUsers } from "@/lib/api";
import type { AuthUser } from "@/lib/types";

export default async function ConversationPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = await params;
  let detail;
  try {
    detail = await getConversation(id);
  } catch (e) {
    if (e instanceof ApiError && e.status === 404) notFound();
    throw e;
  }

  const users = await listUsers().catch((): AuthUser[] => []);

  return (
    <ConversationThread
      conversation={detail.conversation}
      messages={detail.messages}
      users={users}
    />
  );
}
