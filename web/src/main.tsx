import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { createBrowserRouter, RouterProvider } from "react-router";

import { App } from "./App";
import { AuditLog } from "./routes/AuditLog";
import { Channels } from "./routes/Channels";
import { Configuration } from "./routes/Configuration";
import { History } from "./routes/History";
import { Logs } from "./routes/Logs";
import { Members } from "./routes/Members";
import { Messages } from "./routes/Messages";
import { Overview } from "./routes/Overview";
import { Roles } from "./routes/Roles";
import { Settings } from "./routes/Settings";
import { Users } from "./routes/Users";
import { UserProfile } from "./routes/UserProfile";
import "./index.css";

const queryClient = new QueryClient();

const router = createBrowserRouter([
  {
    path: "/",
    element: <App />,
    children: [
      { index: true, element: <Overview /> },
      { path: "logs", element: <Logs /> },
      { path: "history", element: <History /> },
      { path: "discord/members", element: <Members /> },
      { path: "discord/roles", element: <Roles /> },
      { path: "discord/channels", element: <Channels /> },
      { path: "discord/audit-log", element: <AuditLog /> },
      { path: "data/users", element: <Users /> },
      { path: "data/messages", element: <Messages /> },
      { path: "data/settings", element: <Settings /> },
      { path: "users/:id", element: <UserProfile /> },
      { path: "configuration", element: <Configuration /> },
    ],
  },
]);

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <QueryClientProvider client={queryClient}>
      <RouterProvider router={router} />
    </QueryClientProvider>
  </StrictMode>,
);
