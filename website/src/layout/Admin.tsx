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

import { Outlet, useNavigate } from "react-router-dom";
import ThemeToggleButton from "../common/ThemeToggleButton";
import AdminSideNav from "../common/AdminSideNav";
import admin from "../utils/admin";

export default function Admin() {
  const navigate = useNavigate();

  if (!admin.isAuthed()) {
    // navigate("/admin/auth/login");
  }

  return (
    <>
      <main className="flex">
        <AdminSideNav />
        <div className="flex-1">
          <Outlet />
        </div>
      </main>
      <ThemeToggleButton className="fixed bottom-5 right-5" />
    </>
  );
}
