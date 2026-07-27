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

export interface NavItem {
  to: string;
  label: string;
  icon: LucideIcon;
}

export interface NavSection {
  heading?: string;
  items: NavItem[];
}

export const navSections: NavSection[] = [
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

export const navItems: NavItem[] = navSections.flatMap((section) => section.items);
