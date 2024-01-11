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

import DataForm from "../../common/DataForm";
import Icon from "../../common/Icon";
import { Link, useNavigate } from "react-router-dom";
import admin from "../../utils/admin";
import { useState } from "react";
import { twMerge } from "tailwind-merge";

export default function AdminLoginPage() {
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  async function loginHandler(data: Record<string, string>) {
    setLoading(true);
    if (await admin.login(data.username, data.password)) {
      navigate("/admin");
    } else {
      alert("Login failed")
    }
    setLoading(false);
  }

  if(admin.isAuthed()) {
    navigate("/admin")
  }

  return (
    <>
      <section className="h-screen flex justify-center items-center">
        <div className="border border-front border-opacity-20 shadow-xl rounded-xl overflow-hidden">
          <div className="bg-primary flex py-6 px-10 items-center gap-x-3">
            <img
              src="/logo.png"
              className="h-[5em] aspect-square drop-shadow"
            />
            <div className="flex flex-col text-back items-center gap-y-1">
              <h2 className="text-3xl font-medium">arewefastyet</h2>
              <p className="text-sm">admin</p>
            </div>
          </div>
          <DataForm.Container
            className="flex flex-col p-5 gap-y-5"
            onSubmit={(data) => loginHandler(data)}
          >
            <div
              className="relative flex items-center p-2 gap-x-4 border border-front duration-300 rounded-md outline outline-transparent outline-offset-8 bg-foreground bg-opacity-5
          focus-within:-outline-offset-1 focus-within:outline-primary focus-within:border-transparent focus-within:bg-primary focus-within:bg-opacity-5"
            >
              <Icon icon="person" className="text-2xl" />
              <DataForm.Input
                name="username"
                placeholder="Admin username"
                autoComplete="off"
                className="bg-transparent border-none outline-none flex-1"
              />
            </div>

            <div
              className="relative flex items-center p-2 gap-x-4 border border-front duration-300 rounded-md outline outline-transparent outline-offset-8 bg-foreground bg-opacity-5
            focus-within:-outline-offset-1 focus-within:outline-primary focus-within:border-transparent focus-within:bg-primary focus-within:bg-opacity-5"
            >
              <Icon icon="key" className="text-2xl" />
              <DataForm.Input
                name="password"
                placeholder="Admin password"
                autoComplete="off"
                className="bg-transparent border-none outline-none flex-1"
              />
            </div>

            <div className="flex justify-between">
              <Link
                to="/"
                className="px-6 py-2 rounded-md bg-red-600 text-back font-medium"
              >
                Cancel
              </Link>
              <DataForm.Input
                disabled={loading}
                className={twMerge(
                  "px-6 py-2 rounded-md bg-foreground text-back font-medium cursor-pointer",
                  loading && "opacity-50 cursor-not-allowed"
                )}
                type="submit"
                value="Login"
              />
            </div>
          </DataForm.Container>
        </div>
      </section>
    </>
  );
}
