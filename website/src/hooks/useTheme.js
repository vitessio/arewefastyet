import { useState } from "react";
import useGlobalContext from "../contexts/GlobalContext";
import resolveConfig from "tailwindcss/resolveConfig";
import tailwindConfig from "../../tailwind.config.js";

function getThemeColors() {
  const twConf = resolveConfig(tailwindConfig);
  const colors = twConf.theme.colors;

  return {
    primary: colors.primary,
    secondary: colors.secondary,
    accent: colors.accent,
    front: colors.front,
    back: colors.back,
    background: colors.background,
    foreground: colors.foreground,
  };
}

export default function useTheme() {
  const { theme, setTheme } = useGlobalContext();

  const [colors, setColors] = useState(getThemeColors());

  function set(newTheme) {
    setTheme(newTheme);
    setColors(getThemeColors());
  }

  return { current: theme, set, colors };
}
