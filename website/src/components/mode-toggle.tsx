import { Moon, MoonIcon, Sun, SunIcon } from "lucide-react";

import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { useTheme } from "@/components/theme-provider";

export function ModeToggle() {
  const { setTheme } = useTheme();

  const toggleTheme = () => {
    const currentTheme = document.documentElement.classList.contains("dark")
      ? "dark"
      : "light";
    const newTheme = currentTheme === "dark" ? "light" : "dark";
    setTheme(newTheme);
  };

  return (
    <Button
      variant="ghost"
      size={"icon"}
      onClick={toggleTheme}
      className="text-neutral-500"
    >
      <SunIcon className="h-[1.5rem] w-[1.5rem] rotate-0 scale-100 transition-transform dark:-rotate-90 dark:scale-0" />
      <MoonIcon className="absolute h-[1.5rem] w-[1.5rem] rotate-90 scale-0 transition-transform dark:rotate-0 dark:scale-100" />
      <span className="sr-only">Toggle theme</span>
    </Button>
  );
}
