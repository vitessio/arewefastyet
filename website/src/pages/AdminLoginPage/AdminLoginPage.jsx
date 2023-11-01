import React from "react";
import DataForm from "../../common/DataForm";
import Icon from "../../common/Icon";
import { Link } from "react-router-dom";

export default function AdminLoginPage() {
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
            onSubmit={(data) => console.log(data)}
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
              <Link to="/" className="px-6 py-2 rounded-md bg-[#D93036] text-back font-medium"> Cancel</Link>
              <DataForm.Input className="px-6 py-2 rounded-md bg-foreground text-back font-medium cursor-pointer" type="submit" value="Login" />
            </div>
          </DataForm.Container>
        </div>
      </section>
    </>
  );
}
