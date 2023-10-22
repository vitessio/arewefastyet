import useGlobalContext from "../contexts/GlobalContext";

export default function useTheme() {
  const { theme, setTheme } = useGlobalContext();

  function set(newTheme) {
    setTheme(newTheme);
    setColors(getThemeColors());
  }

  return { current: theme, set };
}
