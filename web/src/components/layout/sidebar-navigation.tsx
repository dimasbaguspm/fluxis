import { cx } from "@/lib";
import { getInitials } from "@lib/get-initials";
import { useSessionState } from "@providers/session";
import { vMr2 } from "@versaur/core/utilities";
import { MenuIcon } from "@versaur/icons";
import { Sidebar } from "@versaur/react/blocks";
import { Avatar, ButtonIcon, Heading, Icon } from "@versaur/react/primitive";
import { Link, useLocation } from "react-router";
import { SIDEBAR_NAV_ITEMS } from "./sidebar-nav-config";

export const SidebarNavigation = () => {
  const location = useLocation();
  const { user } = useSessionState();

  const mainItems = SIDEBAR_NAV_ITEMS.filter((item) => item.section === "main");

  const isActive = (href: string) => location.pathname === href;
  const userInitials = getInitials(user?.displayName);

  return (
    <Sidebar>
      <Sidebar.Header>
        <Heading as="h1" weight="bold" size="xl">
          Fluxis
        </Heading>
      </Sidebar.Header>
      <Sidebar.Body>
        <Sidebar.ItemList>
          {mainItems.map((item) => (
            <Sidebar.Item
              key={item.href}
              as={Link}
              to={item.href}
              active={isActive(item.href)}
              icon={<Icon as={item.icon} />}
            >
              {item.label}
            </Sidebar.Item>
          ))}
        </Sidebar.ItemList>
      </Sidebar.Body>
      <Sidebar.Footer>
        <Sidebar.Item
          as="div"
          action={
            <ButtonIcon as={MenuIcon} variant="ghost" size="small" aria-label="User options" />
          }
        >
          <Avatar size="sm" className={cx(vMr2)}>
            {userInitials}
          </Avatar>
          Profile
        </Sidebar.Item>
      </Sidebar.Footer>
    </Sidebar>
  );
};
