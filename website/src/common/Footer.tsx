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

import React from "react";

import { AiFillGithub, AiFillSlackCircle } from "react-icons/ai";
import { FaSquareXTwitter } from "react-icons/fa6";
import { BsStackOverflow } from "react-icons/bs";
import { Link } from "react-router-dom";
import { IconType } from "react-icons";

interface Social {
  url: string;
  icon: IconType;
}

interface LinkItem {
  title: string;
  to: string;
}

interface FooterItem {
  title: string;
  links: LinkItem[];
}

const Footer: React.FC = () => {
  const socials: Social[] = [
    { url: "https://github.com/vitessio/arewefastyet", icon: AiFillGithub },
    { url: "https://vitess.io/slack", icon: AiFillSlackCircle },
    { url: "https://twitter.com/vitessio", icon: FaSquareXTwitter },
    { url: "https://stackoverflow.com/search?q=vitess", icon: BsStackOverflow },
  ];

  const items: FooterItem[] = [
    {
      title: "Benchmarks",
      links: [
        { title: "Daily", to: "/daily" },
        { title: "Compare", to: "/compare" },
        { title: "PR", to: "/pr" },
        { title: "History", to: "/history" },
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
  ];

  return (
    <footer className="m-auto w-full border-border/40 lg:px-60 md:px-12 sm:px-6 py-12 text-foreground border-t">
      <div className="flex flex-col gap-8 lg:gap-28 lg:flex-row md:justify-center text-sm font-light tracking-tight lg:pr-64">
        <div className="md:flex m-auto lg:m-0 flex-col gap-y-4 lg:gap-y-12 lg:order-1 order-2">
          <div className="flex gap-x-2 items-center">
            <img
              src="/logo.png"
              alt="vitess logo"
              className="aspect-square h-10 w-10 md:h-12 md:w-12 md:shrink-0"
            />
            <h2 className="text-3xl font-normal">arewefastyet</h2>
          </div>
          <div className="my-2 flex justify-between md:justify-center px-5 md:px-0 gap-x-4">
            <div className="flex items-center gap-x-8">
              {socials.map((social, key) => (
                <Link key={key} to={social.url} target="__blank" className="">
                  <social.icon className="text-foreground hover:text-primary/80 md:text-4xl text-2xl" />
                </Link>
              ))}
            </div>
          </div>
        </div>
        <div className="flex order-1 md:gap-[12em] gap-[6em] xl:pl-24 m-auto md:justify-center">
          {items.map((item, key) => (
            <div key={key} className="flex flex-col">
              <h5 className="font-bold md:text-lg text-base">{item.title}</h5>
              <div className="my-4 md:my-7 flex flex-col gap-y-2">
                {item.links.map((link, key) => (
                  <Link className="md:text-base text-xs text-foreground/80 hover:text-primary/80" key={key} to={link.to}>
                    {link.title}
                  </Link>
                ))}
              </div>
            </div>
          ))}
        </div>
      </div>
      <div className="mt-8 lg:mt-4 text-center md:text-center text-medium">
        @vitessio/arewefastyet
      </div>
    </footer>
  );
};

export default Footer;