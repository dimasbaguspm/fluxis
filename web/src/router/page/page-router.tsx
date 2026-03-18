import { invalidateQueriesByEvent } from "@/hooks/use-api";
import { Layout } from "@/components/layout";
import { DrawerProvider } from "@/providers/drawer";
import { NotifierProvider } from "@/providers/notifier";
import { SessionHydrator, useSessionState } from "@/providers/session";
import { DrawerRouter } from "@/router/drawer";
import { DEEP_LINKS, PAGES, ROUTE_GROUPS } from "@constants/page-routes";
import { useQueryClient } from "@tanstack/react-query";
import { createBrowserRouter, Navigate, Outlet } from "react-router";
import { BoardsPage } from "./boards-page";
import { DashboardPage } from "./dashboard-page";
import { OrganizationsPage } from "./organizations-page";
import { ProfilePage } from "./profile-page";
import {
  ProjectsPage,
  ProjectDetailLayout,
  ProjectOverviewPage,
  ProjectSprintsPage,
  ProjectBoardPage,
  ProjectTicketsPage,
} from "./projects-page";
import { SettingsPage } from "./settings-page";
import { SignInPage } from "./sign-in-page";
import { SignUpPage } from "./sign-up-page";
import { SprintsPage } from "./sprints-page";
import { TicketsPage } from "./tickets-page";

function CacheInvalidator() {
  const queryClient = useQueryClient();

  return (
    <NotifierProvider
      onNotified={(event) => {
        invalidateQueriesByEvent(queryClient, event);
      }}
    >
      <DrawerProvider>
        <Layout />
        <DrawerRouter />
      </DrawerProvider>
    </NotifierProvider>
  );
}

function ProtectedRoute() {
  const { accessToken } = useSessionState();
  if (!accessToken) {
    return <Navigate to={DEEP_LINKS.SIGN_IN} replace />;
  }

  return (
    <SessionHydrator>
      <CacheInvalidator />
    </SessionHydrator>
  );
}

function UnprotectedRoute() {
  const { accessToken } = useSessionState();
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
            element: <SignInPage />,
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
            children: [
              { index: true, element: <ProjectsPage /> },
              {
                path: ":projectId",
                element: <ProjectDetailLayout />,
                children: [
                  { index: true, element: <ProjectOverviewPage /> },
                  { path: "sprints", element: <ProjectSprintsPage /> },
                  { path: "boards", element: <ProjectBoardPage /> },
                  { path: "tickets", element: <ProjectTicketsPage /> },
                ],
              },
            ],
          },
          {
            path: PAGES.SPRINTS,
            element: <SprintsPage />,
          },
          {
            path: PAGES.BOARDS,
            element: <BoardsPage />,
          },
          {
            path: PAGES.TICKETS,
            element: <TicketsPage />,
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
