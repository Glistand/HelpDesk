import { notFound } from "next/navigation";
import { ConversationThread } from "@/components/conversation-thread";
import { ApiError, getConversation } from "@/lib/api";

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

  return (
    <ConversationThread
      conversation={detail.conversation}
      messages={detail.messages}
    />
  );
}
