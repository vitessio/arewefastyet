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

import React, { createContext, useState, useEffect, useContext } from "react";

const GlobalContext = createContext();

export function GlobalProvider({ children }) {
  const [theme, setTheme] = useState();
  const [modal, setModal] = useState();

  useEffect(() => {
    const storedTheme = localStorage.getItem("vitess__arewefastyet__theme");
    if (storedTheme) {
      setTheme(storedTheme);
    } else {
      const darkColorPreference = window.matchMedia(
        "(prefers-color-scheme: dark)"
      );

      setTheme(darkColorPreference.matches ? "dark" : "default");
    }
  }, []);

  useEffect(() => {
    setTimeout(() => {
      document.documentElement.style.transitionDuration = "700ms";
    }, 100);

    if (theme) {
      document.documentElement.setAttribute("data-theme", theme);
      localStorage.setItem("vitess__arewefastyet__theme", theme);
    }
  }, [theme]);

  const value = { theme, setTheme, modal, setModal };

  return (
    <GlobalContext.Provider value={value}>{children}</GlobalContext.Provider>
  );
}

export default function useGlobalContext() {
  return useContext(GlobalContext);
}
