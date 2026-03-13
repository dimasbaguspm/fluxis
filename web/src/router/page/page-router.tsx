import { createBrowserRouter, Outlet } from "react-router";
import { LoginPage } from "./login-page";

export const router = createBrowserRouter([
  {
    element: <Outlet />,
    children: [
      // unprotected
      {
        path: "/login",
        element: <LoginPage />,
      },
      // protected
      {},
    ],
  },
]);
