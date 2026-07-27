import { useEffect } from "react";
import { NavLink, Outlet } from "react-router";
import { useQuery } from "@tanstack/react-query";
import {
  ClipboardList,
  Hash,
  History as HistoryIcon,
  IdCard,
  LayoutDashboard,
  MessagesSquare,
  ScrollText,
  ShieldCheck,
  SlidersHorizontal,
  Users,
  Wrench,
  type LucideIcon,
} from "lucide-react";

import { UserSearch } from "./components/UserSearch";
import { StatusDot } from "./components/ui/StatusDot";
import { fetchStatus } from "./lib/api";

interface NavItem {
  to: string;
  label: string;
  icon: LucideIcon;
}

interface NavSection {
  heading?: string;
  items: NavItem[];
}

const navSections: NavSection[] = [
  { items: [{ to: "/", label: "Dashboard", icon: LayoutDashboard }] },
  {
    heading: "Operate",
    items: [
      { to: "/logs", label: "Logs", icon: ScrollText },
      { to: "/history", label: "History", icon: HistoryIcon },
    ],
  },
  {
    heading: "Discord",
    items: [
      { to: "/discord/members", label: "Members", icon: Users },
      { to: "/discord/roles", label: "Roles", icon: ShieldCheck },
      { to: "/discord/channels", label: "Channels", icon: Hash },
      { to: "/discord/audit-log", label: "Audit Log", icon: ClipboardList },
    ],
  },
  {
    heading: "Data",
    items: [
      { to: "/data/users", label: "Users", icon: IdCard },
      { to: "/data/messages", label: "Messages", icon: MessagesSquare },
      { to: "/data/settings", label: "Settings", icon: SlidersHorizontal },
    ],
  },
  { items: [{ to: "/configuration", label: "Configuration", icon: Wrench }] },
];

function navLinkClassName({ isActive }: { isActive: boolean }): string {
  return [
    "group flex items-center gap-2.5 rounded-lg px-3 py-1.5 text-sm transition-colors",
    isActive ? "bg-accent-soft text-ink" : "text-ink-muted hover:bg-surface-3 hover:text-ink",
  ].join(" ");
}

function BrandMark() {
  return (
    <div className="flex items-center gap-2.5 px-3">
      <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-gradient-to-br from-accent to-accent-hover text-sm font-bold text-white shadow-[0_0_16px_var(--color-accent-soft)]">
        CJ
      </div>
      <div className="leading-tight">
        <div className="text-sm font-semibold text-ink">CJ Dashboard</div>
        <div className="text-[11px] text-ink-faint">Operator console</div>
      </div>
    </div>
  );
}

function ConnectionFooter() {
  const { data } = useQuery({ queryKey: ["status"], queryFn: fetchStatus, refetchInterval: 5000 });

  return (
    <div className="flex items-center gap-2 px-3 py-2 text-xs text-ink-faint">
      <StatusDot tone={data ? "good" : "neutral"} pulse={!!data} />
      {data ? (
        <span>
          Online · {data.version} · up {data.uptime}
        </span>
      ) : (
        <span>Connecting...</span>
      )}
    </div>
  );
}

function isTypingTarget(target: EventTarget | null): boolean {
  if (!(target instanceof HTMLElement)) return false;
  return target.tagName === "INPUT" || target.tagName === "TEXTAREA" || target.isContentEditable;
}

export function App() {
  useEffect(() => {
    function handleKeyDown(e: KeyboardEvent) {
      if (e.key !== "/" || isTypingTarget(e.target)) return;
      e.preventDefault();
      document.getElementById("global-user-search")?.focus();
    }
    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, []);

  return (
    <div className="flex h-full">
      <nav className="flex w-64 shrink-0 flex-col gap-5 border-r border-border bg-surface p-4">
        <BrandMark />
        <UserSearch />
        <div className="flex flex-1 flex-col gap-5 overflow-y-auto">
          {navSections.map((section, index) => (
            <div key={section.heading ?? index} className="flex flex-col gap-1">
              {section.heading && (
                <div className="px-3 text-[11px] font-semibold uppercase tracking-wider text-ink-faint">
                  {section.heading}
                </div>
              )}
              {section.items.map((item) => {
                const Icon = item.icon;
                return (
                  <NavLink key={item.to} to={item.to} end={item.to === "/"} className={navLinkClassName}>
                    <Icon className="h-4 w-4 shrink-0" strokeWidth={1.75} />
                    {item.label}
                  </NavLink>
                );
              })}
            </div>
          ))}
        </div>
        <div className="border-t border-border pt-2">
          <ConnectionFooter />
        </div>
      </nav>
      <main className="min-w-0 flex-1 overflow-y-auto">
        <Outlet />
      </main>
    </div>
  );
}
