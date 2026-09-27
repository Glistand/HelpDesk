import { redirect } from "next/navigation";
import { getMe, listUsers } from "@/lib/api";
import { initialsFromName } from "@/lib/agents";
import type { AuthUser } from "@/lib/types";

const ROLE_LABEL: Record<string, string> = {
  admin: "Администратор",
  agent: "Агент",
  requester: "Заявитель",
};

export default async function TeamPage() {
  const user = await getMe().catch(() => null);
  if (!user) redirect("/login");
  if (user.role !== "admin") redirect("/inbox");

  const users = await listUsers().catch((): AuthUser[] => []);

  return (
    <div className="mx-auto w-full max-w-2xl px-5 py-8 md:px-8">
      <h1 className="font-ui-semibold text-xl text-primary">Команда</h1>
      <p className="mt-1 text-[13px] text-tertiary">
        Аккаунты агентов и администраторов ({users.length})
      </p>

      {users.length === 0 ? (
        <p className="mt-8 text-[13px] text-tertiary">Нет пользователей</p>
      ) : (
        <ul className="mt-6 divide-y divide-border rounded-md border border-border">
          {users.map((u) => (
            <li key={u.id} className="flex items-center gap-3 px-4 py-3">
              <span className="flex h-9 w-9 shrink-0 items-center justify-center rounded-full bg-accent-muted text-[11px] font-ui-semibold text-accent">
                {initialsFromName(u.name)}
              </span>
              <div className="min-w-0 flex-1">
                <p className="truncate text-[14px] font-ui-medium text-primary">
                  {u.name}
                  {u.id === user.id && (
                    <span className="ml-2 text-[11px] font-normal text-tertiary">вы</span>
                  )}
                </p>
                <p className="truncate text-[12px] text-tertiary">{u.email}</p>
              </div>
              <div className="shrink-0 text-right">
                <p className="text-[12px] text-secondary">
                  {ROLE_LABEL[u.role] || u.role}
                </p>
                <p className="font-mono text-[10px] text-tertiary">{u.id}</p>
              </div>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}
