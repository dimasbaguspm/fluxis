import { ClipboardIcon, GridIcon, HomeIcon, SettingsIcon, UserIcon } from "@versaur/icons";
import { DEEP_LINKS } from "@constants/page-routes";

export interface NavItem {
  label: string;
  href: string;
  icon: typeof HomeIcon;
  section?: "main" | "secondary";
}

export const SIDEBAR_NAV_ITEMS: NavItem[] = [
  {
    label: "Dashboard",
    href: DEEP_LINKS.DASHBOARD,
    icon: HomeIcon,
    section: "main",
  },
  {
    label: "Organizations",
    href: DEEP_LINKS.ORGANIZATIONS,
    icon: GridIcon,
    section: "main",
  },
  {
    label: "Projects",
    href: DEEP_LINKS.PROJECTS,
    icon: ClipboardIcon,
    section: "main",
  },
  {
    label: "Sprints",
    href: DEEP_LINKS.SPRINTS,
    icon: ClipboardIcon,
    section: "main",
  },
  {
    label: "Boards",
    href: DEEP_LINKS.BOARDS,
    icon: GridIcon,
    section: "main",
  },
  {
    label: "Settings",
    href: DEEP_LINKS.SETTINGS,
    icon: SettingsIcon,
    section: "secondary",
  },
  {
    label: "Profile",
    href: DEEP_LINKS.PROFILE,
    icon: UserIcon,
    section: "secondary",
  },
];
