import React from "react";
import { twMerge } from "tailwind-merge";
import useTheme from "../hooks/useTheme";
import Icon from "./Icon";

export default function ThemeToggleButton(props: { className?: string }) {
  const theme = useTheme();

  function toggleTheme() {
    theme.set((t) => (t === "default" ? "dark" : "default"));
  }

  return (
    <button
      className={twMerge(
        "bg-foreground p-5 rounded-full text-back",
        props.className
      )}
      onClick={toggleTheme}
    >
      <Icon icon={(theme.current === "dark" && "light_mode") || "dark_mode"} />
    </button>
  );
}
