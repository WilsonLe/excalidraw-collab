import { Excalidraw, serializeAsJSON } from "@excalidraw/excalidraw";
import { ExcalidrawElement } from "@excalidraw/excalidraw/types/element/types";
import {
  AppState,
  BinaryFiles,
  ExcalidrawImperativeAPI,
} from "@excalidraw/excalidraw/types/types";
import React, { useEffect, useRef, useState } from "react";
import "./App.css";

function App() {
  const [excalidrawAPI, setExcalidrawAPI] =
    useState<ExcalidrawImperativeAPI | null>(null);
  const [ws, setWs] = useState<WebSocket | null>(null);
  const [sessionId, setSessionId] = useState<string | null>(null);

  useEffect(() => {
    const ws = new WebSocket("/ws");
    ws.addEventListener("open", () => {
      console.log("WebSocket connection established");
    });
    ws.addEventListener("close", () => {
      console.log("WebSocket connection closed");
    });
    ws.addEventListener("error", (error) => {
      console.error("WebSocket error: ", error);
    });
    setWs(ws);
  }, []);

  useEffect(() => {
    if (!ws) return;
    ws.addEventListener(
      "message",
      (event) => {
        if (
          typeof event.data === "string" &&
          event.data.startsWith("session_id_")
        ) {
          setSessionId(event.data.replace("session_id_", ""));
        }
      },
      { once: true },
    );
  }, [ws]);

  useEffect(() => {
    if (!ws || !sessionId) return;
    ws.addEventListener("message", (event) => {
      const data = JSON.parse(event.data);
      const incomingDataSessionId = data.sessionId;
      if (typeof incomingDataSessionId !== "string") return;
      if (incomingDataSessionId === sessionId) return;
      const excalidrawData = JSON.parse(data.excalidrawData);
      excalidrawAPI?.updateScene({
        appState: excalidrawData.appState,
        elements: excalidrawData.elements,
      });
    });
  }, [ws, sessionId]);

  const handleChange = (
    elements: readonly ExcalidrawElement[],
    appState: AppState,
    files: BinaryFiles,
  ) => {
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(serializeAsJSON(elements, appState, files, "local"));
    }
  };

  return (
    <div style={{ height: "100vh", width: "100vw" }}>
      <Excalidraw
        isCollaborating={true}
        excalidrawAPI={setExcalidrawAPI}
        onChange={(elements, appState, files) =>
          handleChange(elements, appState, files)
        }
      />
    </div>
  );
}

export default App;
