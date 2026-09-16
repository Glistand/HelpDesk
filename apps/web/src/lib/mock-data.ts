import type { Agent, Ticket } from "./types";

export const currentAgent: Agent = {
  id: "a-1",
  name: "Алексей К.",
  initials: "АК",
  role: "L1 Support",
};

export const agents: Agent[] = [
  currentAgent,
  { id: "a-2", name: "Марина С.", initials: "МС", role: "L1 Support" },
  { id: "a-3", name: "Денис В.", initials: "ДВ", role: "L2 Escalation" },
];

export const tickets: Ticket[] = [
  {
    id: "t-1042",
    title: "VPN не подключается после обновления клиента",
    description:
      "После обновления Cisco AnyConnect до 5.x соединение обрывается через 30 секунд. Ошибка: «Connection attempt failed». Пробовал перезагрузку ноутбука — не помогло.",
    status: "open",
    priority: "high",
    category: "Сеть / VPN",
    requester: "Иван Петров · Отдел продаж",
    assignee: currentAgent,
    createdAt: "2026-09-16T06:12:00Z",
    updatedAt: "2026-09-16T07:05:00Z",
    sla: {
      firstResponseDue: "2026-09-16T06:42:00Z",
      resolveDue: "2026-09-16T10:12:00Z",
      state: "warning",
    },
    timeline: [
      {
        id: "e-1",
        kind: "created",
        title: "Тикет создан",
        detail: "Канал: портал самообслуживания",
        actor: "Иван Петров",
        at: "2026-09-16T06:12:00Z",
      },
      {
        id: "e-2",
        kind: "assigned",
        title: "Назначен агент",
        detail: "Round-robin · очередь L1",
        actor: "assignment-service",
        at: "2026-09-16T06:12:04Z",
      },
      {
        id: "e-3",
        kind: "notification",
        title: "Уведомление отправлено",
        detail: "Email агенту · mock",
        at: "2026-09-16T06:12:05Z",
      },
      {
        id: "e-4",
        kind: "sla_warn",
        title: "SLA: первый ответ",
        detail: "Осталось 8 минут до дедлайна",
        at: "2026-09-16T06:34:00Z",
      },
    ],
  },
  {
    id: "t-1041",
    title: "Не работает почта в Outlook",
    description: "Outlook показывает «Disconnected». Webmail открывается нормально.",
    status: "pending",
    priority: "normal",
    category: "Почта",
    requester: "Ольга Смирнова · HR",
    assignee: agents[1],
    createdAt: "2026-09-16T05:40:00Z",
    updatedAt: "2026-09-16T06:55:00Z",
    sla: {
      firstResponseDue: "2026-09-16T06:10:00Z",
      resolveDue: "2026-09-16T09:40:00Z",
      state: "ok",
    },
    timeline: [
      {
        id: "e-1",
        kind: "created",
        title: "Тикет создан",
        actor: "Ольга Смирнова",
        at: "2026-09-16T05:40:00Z",
      },
      {
        id: "e-2",
        kind: "assigned",
        title: "Назначен агент",
        actor: "assignment-service",
        at: "2026-09-16T05:40:03Z",
      },
      {
        id: "e-3",
        kind: "comment",
        title: "Комментарий агента",
        detail: "Попробуйте создать новый профиль Outlook. Инструкция отправлена.",
        actor: "Марина С.",
        at: "2026-09-16T06:02:00Z",
      },
      {
        id: "e-4",
        kind: "status",
        title: "Статус изменён",
        detail: "open → pending (ожидание ответа пользователя)",
        actor: "Марина С.",
        at: "2026-09-16T06:55:00Z",
      },
    ],
  },
  {
    id: "t-1038",
    title: "Запрос доступа к CRM для нового сотрудника",
    description: "Нужен доступ read-only к Bitrix24 для стажёра отдела маркетинга.",
    status: "new",
    priority: "low",
    category: "Доступы",
    requester: "Екатерина Л. · Маркетинг",
    createdAt: "2026-09-16T07:08:00Z",
    updatedAt: "2026-09-16T07:08:00Z",
    sla: {
      firstResponseDue: "2026-09-16T07:38:00Z",
      resolveDue: "2026-09-16T11:08:00Z",
      state: "ok",
    },
    timeline: [
      {
        id: "e-1",
        kind: "created",
        title: "Тикет создан",
        actor: "Екатерина Л.",
        at: "2026-09-16T07:08:00Z",
      },
    ],
  },
  {
    id: "t-1035",
    title: "Принтер в open-space не печатает",
    description: "HP LaserJet на 3 этаже — задания висят в очереди, ошибок на дисплее нет.",
    status: "open",
    priority: "urgent",
    category: "Оборудование",
    requester: "Антон М. · Офис-менеджер",
    assignee: agents[2],
    createdAt: "2026-09-15T14:20:00Z",
    updatedAt: "2026-09-16T06:30:00Z",
    sla: {
      firstResponseDue: "2026-09-15T14:50:00Z",
      resolveDue: "2026-09-15T18:20:00Z",
      state: "breached",
    },
    timeline: [
      {
        id: "e-1",
        kind: "created",
        title: "Тикет создан",
        at: "2026-09-15T14:20:00Z",
      },
      {
        id: "e-2",
        kind: "assigned",
        title: "Назначен агент",
        at: "2026-09-15T14:20:02Z",
      },
      {
        id: "e-3",
        kind: "sla_breach",
        title: "SLA нарушен",
        detail: "Resolve time exceeded",
        at: "2026-09-15T18:20:00Z",
      },
      {
        id: "e-4",
        kind: "escalated",
        title: "Эскалация на L2",
        detail: "Приоритет повышен · очередь escalation",
        at: "2026-09-15T18:20:05Z",
      },
      {
        id: "e-5",
        kind: "notification",
        title: "Уведомление L2",
        detail: "Webhook mock → escalation channel",
        at: "2026-09-15T18:20:06Z",
      },
    ],
  },
  {
    id: "t-1030",
    title: "Сброс пароля Active Directory",
    description: "Пользователь заблокирован после 5 неудачных попыток входа.",
    status: "resolved",
    priority: "normal",
    category: "Доступы",
    requester: "Сергей К. · Бухгалтерия",
    assignee: currentAgent,
    createdAt: "2026-09-15T09:15:00Z",
    updatedAt: "2026-09-15T09:28:00Z",
    sla: {
      firstResponseDue: "2026-09-15T09:45:00Z",
      resolveDue: "2026-09-15T13:15:00Z",
      state: "ok",
    },
    timeline: [
      {
        id: "e-1",
        kind: "created",
        title: "Тикет создан",
        at: "2026-09-15T09:15:00Z",
      },
      {
        id: "e-2",
        kind: "assigned",
        title: "Назначен агент",
        at: "2026-09-15T09:15:01Z",
      },
      {
        id: "e-3",
        kind: "status",
        title: "Статус изменён",
        detail: "open → resolved",
        actor: "Алексей К.",
        at: "2026-09-15T09:28:00Z",
      },
    ],
  },
];

export function getTicket(id: string): Ticket | undefined {
  return tickets.find((t) => t.id === id);
}

export const inboxStats = {
  open: tickets.filter((t) => ["new", "open", "pending"].includes(t.status)).length,
  breached: tickets.filter((t) => t.sla.state === "breached").length,
  mine: tickets.filter((t) => t.assignee?.id === currentAgent.id).length,
};
