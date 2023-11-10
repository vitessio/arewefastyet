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

import {
  AiFillGithub,
  AiFillSlackCircle,
  AiFillTwitterCircle,
} from "react-icons/ai";
import { BsStackOverflow } from "react-icons/bs";
import { Link } from "react-router-dom";

const Footer = () => {
  const socials = [
    { url: "https://github.com/vitessio/arewefastyet", icon: AiFillGithub },
    { url: "https://vitess.io/slack", icon: AiFillSlackCircle },
    { url: "https://twitter.com/vitessio", icon: AiFillTwitterCircle },
    { url: "https://stackoverflow.com/search?q=vitess", icon: BsStackOverflow },
  ];

  const items = [
    {
      title: "Benchmarks",
      links: [
        { title: "Daily", to: "/daily" },
        { title: "Compare", to: "/compare" },
        { title: "Micro", to: "/micro" },
        { title: "Macro", to: "/macro" },
        { title: "PR", to: "/pr" },
      ],
    },
    {
      title: "Vitess",
      links: [
        { title: "Docs", to: "https://vitess.io/docs" },
        { title: "Blog", to: "https://vitess.io/blog" },
        { title: "Community", to: "https://vitess.io/community" },
      ],
    },
    {
      title: "arewefastyet",
      links: [{ title: "GitHub", to: "https://www.github.com/vitessio/arewefastyet" }],
    },
  ];

  return (
    <footer className="p-page relative mt-10 py-20 text-back">
      <div className="absolute bottom-6 left-6 right-6 top-6 -z-1 rounded-lg bg-primary" />
      <div className="flex justify-between text-sm font-light tracking-tight">
        <div
          className={`flex basis-[20%] flex-col items-center gap-y-5 opacity-90 drop-shadow-[0px_0px_0px_theme("colors.back")]`}
        >
          <img
            src="/logo.png"
            alt="agrosurance logo"
            className="aspect-square w-1/3"
          />
          <h2 className="text-xl font-medium">arewefastyet</h2>
        </div>
        {items.map((item, key) => (
          <div key={key} className="flex flex-col">
            <h5 className="font-semibold">{item.title}</h5>
            <div className="my-7 flex flex-col gap-y-3">
              {item.links.map((link, key) => (
                <Link key={key} to={link.to}>{link.title}</Link>
              ))}
            </div>
          </div>
        ))}

        {/* <div className="flex flex-col items-center text-center text-back text-opacity-80">
          <h5 className="my-3 text-4xl font-bold tracking-tighter text-back">
            TRULY TESTED
          </h5>
          <p>A cutting edge Benchmarking</p>
          <p>approach for uparalleled</p>
          <p>Database Speed</p>
          <p className="font-semibold text-back">for vitess</p>
        </div> */}
        <div className="w-[15%]" />
      </div>
      <div className="my-2 flex gap-x-4">
        <p className="font-mono text-xl font-bold">Follow Us</p>
        <div className="flex items-center gap-x-3">
          {socials.map((social, key) => (
            <Link key={key} to={social.url} target="__blank" className="">
              <social.icon size={33} />
            </Link>
          ))}
        </div>
      </div>
      <div className="my-2 w-full border border-back"></div>
      <div className="mt-4 text-xs">@vitessio/arewefastyet</div>
    </footer>
  );
};

export default Footer;
