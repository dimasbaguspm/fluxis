import { BottomBar } from "@versaur/react/blocks";
import { Icon } from "@versaur/react/primitive";
import { useLocation, useNavigate } from "react-router";
import { BOTTOM_BAR_NAV_ITEMS } from "./mobile-nav-config";

export const BottomBarNavigation = () => {
  const location = useLocation();
  const navigate = useNavigate();

  const isActive = (href?: string) => {
    if (!href) return false;
    return location.pathname === href;
  };

  const handleNavigate = (href?: string) => {
    if (href) {
      navigate(href);
    }
  };

  return (
    <BottomBar>
      {BOTTOM_BAR_NAV_ITEMS.map((item) => {
        if (item.isMore) {
          return (
            <BottomBar.Item key={item.label} icon={<Icon as={item.icon} size="md" />}>
              {item.label}
            </BottomBar.Item>
          );
        }

        return (
          <BottomBar.Item
            key={item.href}
            onClick={() => handleNavigate(item.href)}
            active={isActive(item.href)}
            icon={<Icon as={item.icon} size="md" />}
          >
            {item.label}
          </BottomBar.Item>
        );
      })}
    </BottomBar>
  );
};
