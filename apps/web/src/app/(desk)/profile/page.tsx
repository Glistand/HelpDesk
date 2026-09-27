import Link from "next/link";
import { redirect } from "next/navigation";
import { getMe } from "@/lib/api";
import { initialsFromName } from "@/lib/agents";

export default async function ProfilePage() {
  const user = await getMe().catch(() => null);
  if (!user) redirect("/login");

  return (
    <div className="mx-auto w-full max-w-lg px-5 py-8 md:px-8">
      <h1 className="font-ui-semibold text-xl text-primary">Профиль</h1>
      <p className="mt-1 text-[13px] text-tertiary">
        Данные аккаунта. Смена пароля пока недоступна.
      </p>

      <div className="mt-8 flex items-center gap-4">
        <span className="flex h-14 w-14 items-center justify-center rounded-full bg-accent-muted text-lg font-ui-semibold text-accent">
          {initialsFromName(user.name)}
        </span>
        <div>
          <p className="font-ui-semibold text-primary">{user.name}</p>
          <p className="text-[13px] text-tertiary">{user.email}</p>
        </div>
      </div>

      <dl className="mt-8 space-y-4 border-t border-border pt-6">
        <div>
          <dt className="text-[11px] uppercase tracking-wide text-tertiary">ID</dt>
          <dd className="mt-0.5 font-mono text-[13px] text-secondary">{user.id}</dd>
        </div>
        <div>
          <dt className="text-[11px] uppercase tracking-wide text-tertiary">Роль</dt>
          <dd className="mt-0.5 text-[13px] text-secondary">{user.role}</dd>
        </div>
        <div>
          <dt className="text-[11px] uppercase tracking-wide text-tertiary">Email</dt>
          <dd className="mt-0.5 text-[13px] text-secondary">{user.email}</dd>
        </div>
      </dl>

      {user.role === "admin" && (
        <p className="mt-8 text-[13px] text-tertiary">
          Управление аккаунтами — в{" "}
          <Link href="/team" className="text-accent underline">
            Команде
          </Link>
          .
        </p>
      )}
    </div>
  );
}
