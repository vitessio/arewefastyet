/*
Copyright 2024 The Vitess Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import { Link, NavLink } from "react-router-dom";
import { twMerge } from "tailwind-merge";
import { ModeToggle } from "@/components/mode-toggle";
import { Button } from "@/components/ui/button";
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from "@/components/ui/sheet";

const navItems = [
  { to: "/status", title: "Status" },
  { to: "/daily", title: "Daily" },
  { to: "/compare", title: "Compare" },
  { to: "/search", title: "Search" },
  { to: "/macro", title: "Macro" },
  { to: "/fk", title: "Foreign Keys" },
  { to: "/pr", title: "PR" },
];

export default function Navbar() {
  return (
    <nav className="flex flex-col relative">
      <div
        className={twMerge(
          "w-full bg-background z-[49] flex justify-between md:justify-center p-page py-4 border-b border-border"
        )}
      >
        <Link to="/" className="flex flex-1 gap-x-2 items-center">
          <img src="/logo.png" className="h-[2em]" alt="logo" />
          <h1 className="font-semilight text-lg md:hidden lg:block">arewefastyet</h1>
        </Link>

        <div className="hidden md:flex gap-x-10 items-center">
          {navItems.map((item, key) => (
            <NavLink
              key={key}
              to={item.to}
              className={({ isActive, isPending }) =>
                twMerge(
                  "text-lg",
                  isPending
                    ? "pointer-events-none opacity-50"
                    : isActive
                    ? "text-primary"
                    : ""
                )
              }
            >
              {item.title}
            </NavLink>
          ))}
        </div>

        <div className="flex-1 flex gap-3 justify-end items-center">
          <ModeToggle />

          <Sheet>
            <SheetTrigger asChild>
              <Button
                variant={"ghost"}
                size={"icon"}
                className="relative md:hidden text-3xl flex items-center"
              >
                <span className="sr-only">Open Menu</span>
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="24"
                  height="24"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  className="lucide lucide-menu"
                >
                  <line x1="4" x2="20" y1="12" y2="12" />
                  <line x1="4" x2="20" y1="6" y2="6" />
                  <line x1="4" x2="20" y1="18" y2="18" />
                </svg>
              </Button>
            </SheetTrigger>
            <SheetContent className="border-accent">
              <SheetHeader>
                <SheetTitle className="text-left">
                  Vitess | Arewefastyet
                </SheetTitle>
                <SheetDescription asChild>
                  <div className="flex flex-col justify-center items-center z-[50] bg-background w-full">
                    <NavLink
                      to={"/"}
                      className={({ isActive, isPending }) =>
                        twMerge(
                          "text-xl text-left font-medium py-3 w-full",
                          isPending
                            ? "pointer-events-none "
                            : isActive
                            ? "text-foreground"
                            : ""
                        )
                      }
                    >
                      Home
                    </NavLink>
                    {navItems.map((item, key) => (
                      <NavLink
                        key={key}
                        to={item.to}
                        className={({ isActive, isPending }) =>
                          twMerge(
                            "text-xl text-left font-medium py-3 w-full",
                            isPending
                              ? "pointer-events-none "
                              : isActive
                              ? "text-foreground"
                              : ""
                          )
                        }
                      >
                        {item.title}
                      </NavLink>
                    ))}
                  </div>
                </SheetDescription>
              </SheetHeader>
            </SheetContent>
          </Sheet>
        </div>
      </div>
    </nav>
  );
}
