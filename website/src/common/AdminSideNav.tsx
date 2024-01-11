import React, { useEffect, useRef, useState } from "react";
import { NavLink } from "react-router-dom";
import { twMerge } from "tailwind-merge";
import Icon from "./Icon";
import admin from "../utils/admin";

// path relative to "/admin"
const navItems = [
  { title: "Dashboard", to: "/", icon: "person" },
  { title: "Benchmarks", to: "/benchmarks", icon: "ssidChart" },
  { title: "Logs", to: "/logs", icon: "description" },
  // { title: "Dummy", to: "/dummy", icon: "key" },
] as const;

export default function AdminSideNav() {
  const [width, setWidth] = useState(0);

  const navRef = useRef() as React.MutableRefObject<HTMLElement>;

  useEffect(() => setWidth(navRef.current.offsetWidth), []);

  return (
    <>
      <div style={{ width }} />
      <nav
        ref={navRef}
        className="h-screen flex flex-col p-6 bg-foreground bg-opacity-10 fixed"
      >
        <div className="flex items-center gap-x-1 text-3xl pl-3">
          <img
            src="/logo.png"
            alt="vitess logo"
            className="h-[1.23em] aspect-square"
          />
          <h3 className="font-medium tracking-tight text-primary font-opensans relative after:absolute after:top-full after:right-1 after:content-['admin'] after:text-xs">
            arewefastyet
          </h3>
        </div>

        <div className="flex-1 flex flex-col gap-y-1 my-8">
          {navItems.map((item, key) => (
            <AdminNavLink
              to={`/admin${item.to}`}
              key={key}
              className="flex gap-x-3 items-center"
              icon={item.icon}
            >
              {item.title}
            </AdminNavLink>
          ))}
        </div>

        <button
          onClick={() => admin.logout()}
          className="bg-red-600 text-white shadow w-full px-5 py-2 rounded-xl duration-150 flex items-center gap-x-3"
        >
          <Icon icon={"logout"} className="text-2xl" /> Logout
        </button>
      </nav>
    </>
  );
}

interface AdminNavLinkProps {
  children?: React.ReactNode;
  className?: string;
  icon?: React.ComponentPropsWithoutRef<typeof Icon>["icon"];
  to: string;
}

function AdminNavLink(props: AdminNavLinkProps) {
  return (
    <NavLink
      className={({ isActive, isPending }) =>
        twMerge(
          "w-full px-5 py-2 rounded-xl duration-150 flex items-center gap-x-3",
          isPending
            ? "pointer-events-none opacity-50"
            : isActive
            ? "bg-primary pointer-events-none"
            : "hover:bg-background",
          props.className
        )
      }
      to={props.to}
    >
      {props.icon && <Icon icon={props.icon} className="text-2xl" />}
      {props.children}
    </NavLink>
  );
}
