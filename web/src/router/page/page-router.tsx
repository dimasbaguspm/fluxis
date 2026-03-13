import { PAGES, ROUTE_GROUPS } from "@constants/page-routes";
import { createBrowserRouter, Outlet } from "react-router";
import { BoardsPage } from "./boards-page";
import { DashboardPage } from "./dashboard-page";
import { LoginPage } from "./login-page";
import { OrganizationsPage } from "./organizations-page";
import { ProfilePage } from "./profile-page";
import { ProjectsPage } from "./projects-page";
import { SettingsPage } from "./settings-page";
import { SignUpPage } from "./sign-up-page";

export const router = createBrowserRouter([
  {
    element: <Outlet />,
    children: [
      // unprotected
      {
        id: ROUTE_GROUPS.UNPROTECTED,
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
      // protected
      {
        id: ROUTE_GROUPS.PROTECTED,
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
