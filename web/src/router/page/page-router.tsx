import { DEEP_LINKS, PAGES, ROUTE_GROUPS } from "@constants/page-routes";
import { useSessionStore } from "@providers/session";
import { createBrowserRouter, Navigate, Outlet } from "react-router";
import { BoardsPage } from "./boards-page";
import { DashboardPage } from "./dashboard-page";
import { LoginPage } from "./login-page";
import { OrganizationsPage } from "./organizations-page";
import { ProfilePage } from "./profile-page";
import { ProjectsPage } from "./projects-page";
import { SettingsPage } from "./settings-page";
import { SignUpPage } from "./sign-up-page";

function ProtectedRoute() {
  const accessToken = useSessionStore((state) => state.accessToken);
  if (!accessToken) {
    return <Navigate to={DEEP_LINKS.SIGN_IN} replace />;
  }
  return <Outlet />;
}

function UnprotectedRoute() {
  const accessToken = useSessionStore((state) => state.accessToken);
  if (accessToken) {
    return <Navigate to={DEEP_LINKS.DASHBOARD} replace />;
  }
  return <Outlet />;
}

export const router = createBrowserRouter([
  {
    element: <Outlet />,
    children: [
      {
        id: ROUTE_GROUPS.UNPROTECTED,
        element: <UnprotectedRoute />,
        children: [
          {
            path: PAGES.SIGN_IN,
            element: <LoginPage />,
          },
          {
            path: PAGES.SIGN_UP,
            element: <SignUpPage />,
          },
        ],
      },
      {
        id: ROUTE_GROUPS.PROTECTED,
        element: <ProtectedRoute />,
        children: [
          {
            path: PAGES.DASHBOARD,
            element: <DashboardPage />,
          },
          {
            path: PAGES.ORGANIZATIONS,
            element: <OrganizationsPage />,
          },
          {
            path: PAGES.PROJECTS,
            element: <ProjectsPage />,
          },
          {
            path: PAGES.BOARDS,
            element: <BoardsPage />,
          },
          {
            path: PAGES.SETTINGS,
            element: <SettingsPage />,
          },
          {
            path: PAGES.PROFILE,
            element: <ProfilePage />,
          },
        ],
      },
    ],
  },
]);
