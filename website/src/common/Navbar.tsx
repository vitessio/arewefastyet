/*
Copyright 2023 The Vitess Authors.

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

import React, { useState } from "react";
import { Link, NavLink } from "react-router-dom";
import { twMerge } from "tailwind-merge";
import Icon from "./Icon";
import useTheme from "../hooks/useTheme";

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
  const [isMenuOpen, setIsMenuOpen] = useState(false);

  const theme = useTheme();

  function toggleMenu() {
    setIsMenuOpen((o) => !o);
  }

  function toggleTheme() {
    theme.set((t) => (t === "default" ? "dark" : "default"));
  }

  return (
    <nav className="flex flex-col relative">
      <div
        className={twMerge(
          "w-full bg-background z-[999] flex justify-between md:justify-center p-page py-4 border-b border-front border-opacity-30 duration-500"
        )}
      >
        <Link to="/" className="flex flex-1 gap-x-2 items-center">
          <img src="/logo.png" className="h-[2.5em]" alt="logo" />
          <h1 className="hidden lg:block font-medium text-lg md:text-2xl">
            arewefastyet
          </h1>
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
          <button
            className="relative text-3xl flex items-center"
            onClick={toggleTheme}
          >
            <Icon
              icon={(theme.current === "dark" && "light_mode") || "dark_mode"}
            />
          </button>
          <button
            className="relative md:hidden text-3xl flex items-center"
            onClick={toggleMenu}
          >
            {isMenuOpen ? (
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
                className="lucide lucide-x"
              >
                <path d="M18 6 6 18" />
                <path d="m6 6 12 12" />
              </svg>
            ) : (
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
            )}
          </button>
        </div>
      </div>
      {isMenuOpen && (
        <div className="flex flex-col justify-center items-center z-[999] bg-background w-full">
          {navItems.map((item, key) => (
            <NavLink
              key={key}
              to={item.to}
              className={({ isActive, isPending }) =>
                twMerge(
                  "text-2xl text-center font-medium border-b border-front border-opacity-30 py-3 w-full",
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
      )}
    </nav>
  );
}
