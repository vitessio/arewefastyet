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

import CompareActionsWrapper from "@/common/CompareActionsWrapper";
import Footer from "@/common/Footer";
import Navbar from "@/common/Navbar";
import { ThemeProvider } from "@/components/theme-provider";
import { CompareProvider } from "@/contexts/CompareContext";
import { Outlet } from "react-router-dom";

export default function Layout() {
  return (
    <ThemeProvider defaultTheme="system" storageKey="vite-ui-theme">
      <CompareProvider>
        <div className="flex flex-col min-h-screen">
          <Navbar />
          <div className="flex-1">
            <Outlet />
            <CompareActionsWrapper />
          </div>
          <Footer />
        </div>
      </CompareProvider>
    </ThemeProvider>
  );
}
