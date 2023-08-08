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

import "./footer.css";
import logo from "../../assets/logo.png";
import {
  AiFillGithub,
  AiFillSlackCircle,
  AiFillTwitterCircle,
} from "react-icons/ai";
import { BsStackOverflow } from "react-icons/bs";

const Footer = () => {
  const icons = [
    { url: "https://github.com/vitessio/arewefastyet", icon: AiFillGithub },
    { url: "https://vitess.io/slack", icon: AiFillSlackCircle },
    { url: "https://twitter.com/vitessio", icon: AiFillTwitterCircle },
    {
      url: "https://stackoverflow.com/search?q=vitess",
      icon: BsStackOverflow,
    },
  ];

  return (
    <div className="footer__main">
      <div className="footer">
        <div className="footer__one">
          <div className="footer__company">
            <img src={logo} alt="logo" className="footer__logo" />
            <h1>Arewefastyet</h1>
          </div>
          <div className="footer__icons">
            {icons.map((item, index) => (
              <a key={index} href={item.url} target="_blank">
                <item.icon size={33} />
              </a>
            ))}
          </div>
        </div>
        <div className="links">
          <div className="footer__usefulLinks">
            <h2>Vitess</h2>
            <div className="">
              <a href="https://vitess.io/docs" target="_blank">
                <p>Docs</p>
              </a>
              <a href="https://vitess.io/blog" target="_blank">
                <p>Blog</p>
              </a>
              <a href="https://vitess.io/community" target="_blank">
                <p>Community</p>
              </a>
            </div>
          </div>
          <div className="footer__usefulLinks">
            <a href="https://github.com/vitessio/arewefastyet" target="_blank">
              <h2>Arewefastyet</h2>
            </a>
            <div className="">
              <a
                href="https://github.com/vitessio/arewefastyet"
                target="_blank"
              >
                <p>GitHub</p>
              </a>
            </div>
          </div>
        </div>
      </div>
      <div>
        <div className="footer__line"></div>
        <p className="copyrightText__footer justify--content">
          Copyright © 2023
        </p>
      </div>
    </div>
  );
};

export default Footer;
