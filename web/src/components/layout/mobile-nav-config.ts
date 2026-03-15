import { DEEP_LINKS } from "@constants/page-routes";
import { BellIcon, HomeIcon, MenuIcon, SearchIcon, SettingsIcon } from "@versaur/icons";

export interface MobileNavItem {
  label: string;
  href?: string;
  icon: typeof HomeIcon;
  isMore?: boolean;
}

export const BOTTOM_BAR_NAV_ITEMS: MobileNavItem[] = [
  {
    label: "Dashboard",
    href: DEEP_LINKS.DASHBOARD,
    icon: HomeIcon,
  },
  {
    label: "Organizations",
    href: DEEP_LINKS.ORGANIZATIONS,
    icon: SearchIcon,
  },
  {
    label: "Projects",
    href: DEEP_LINKS.PROJECTS,
    icon: BellIcon,
  },
  {
    label: "Settings",
    href: DEEP_LINKS.SETTINGS,
    icon: SettingsIcon,
  },
  {
    label: "More",
    icon: MenuIcon,
    isMore: true,
  },
];
