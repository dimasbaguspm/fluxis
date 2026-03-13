import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "@versaur/core/tokens.css";
import { Button } from "@versaur/react/primitive";

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <Button>Hehe</Button>
  </StrictMode>,
);
