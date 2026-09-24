import { getProject } from "@/lib/api";

export default async function ProjectPage() {
  const project = await getProject();
  return (
    <div className="mx-auto w-full max-w-3xl p-8">
      <p className="text-xs uppercase tracking-wide text-tertiary">Проект</p>
      <h1 className="mt-2 text-2xl font-ui-semibold text-primary">{project.name}</h1>
      <p className="mt-3 text-sm text-secondary">{project.description}</p>
      <dl className="mt-8 grid gap-4 rounded-lg border border-border p-5 text-sm sm:grid-cols-2">
        <div><dt className="text-tertiary">Домен</dt><dd className="mt-1 text-primary">{project.domain || "—"}</dd></div>
        <div><dt className="text-tertiary">Часовой пояс</dt><dd className="mt-1 text-primary">{project.timezone}</dd></div>
        <div><dt className="text-tertiary">Email поддержки</dt><dd className="mt-1 text-primary">{project.supportEmail || "—"}</dd></div>
        <div><dt className="text-tertiary">Ключ виджета</dt><dd className="mt-1 font-mono text-primary">{project.siteKey}</dd></div>
      </dl>
    </div>
  );
}
