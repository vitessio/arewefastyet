import React from "react";
import { NavLink } from "react-router-dom";
import { twMerge } from "tailwind-merge";

// path relative to "/admin"
const navItems = [
  { title: "Dashboard", to: "/" },
  { title: "Benchmarks", to: "/benchmarks" },
  { title: "Logs", to: "/logs" },
  { title: "Dummy", to: "/dummy" },
];

export default function AdminSideNav() {
  return (
    <nav className="h-screen flex flex-col p-6 bg-foreground bg-opacity-10">
      <div className="flex items-center gap-x-1 text-3xl pl-3">
        <img
          src="/logo.png"
          alt="vitess logo"
          className="h-[1.23em] aspect-square"
        />
        <h3 className="font-medium tracking-tight text-primary font-opensans">
          arewefastyet
        </h3>
      </div>

      <div className="flex-1 flex flex-col gap-y-1 my-8">
        {navItems.map((item, key) => (
          <AdminNavLink to={`/admin${item.to}`} key={key}>
            {item.title}
          </AdminNavLink>
        ))}
      </div>

      <AdminNavLink
        to="/admin/auth/logout"
        className="bg-red-600 text-white shadow"
      >
        Logout
      </AdminNavLink>
    </nav>
  );
}

interface AdminNavLinkProps {
  children?: React.ReactNode;
  className?: string;
  to: string;
}

function AdminNavLink(props: AdminNavLinkProps) {
  return (
    <NavLink
      className={({ isActive, isPending }) =>
        twMerge(
          "w-full px-5 py-2 rounded-xl duration-150",
          isPending
            ? "pointer-events-none opacity-50"
            : isActive
            ? "bg-primary"
            : "hover:bg-background",
          props.className
        )
      }
      to={props.to}
    >
      {props.children}
    </NavLink>
  );
}
