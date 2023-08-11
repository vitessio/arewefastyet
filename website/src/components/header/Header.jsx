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

import React, { useState, useEffect, useContext } from "react";
import { NavLink } from "react-router-dom";

import "./header.css";

import { AppContext } from "../../AppContext";
import logo from "../../assets/logo.png";

const Header = () => {
  const { isColorChanged, handleButtonClick } = useContext(AppContext);

  const [openMenu, setOpenMenu] = useState(false);
  const handleMenuClick = () => {
    window.scrollTo(0, 0);
    if (window.innerWidth <= 768) {
      setOpenMenu(false);
    } else {
      setOpenMenu(true);
    }
  };

  const [visible, setVisible] = useState(true);
  const [lastScrollTop, setLastScrollTop] = useState(0);

  useEffect(() => {
    const handleScroll = () => {
      const scrollTop =
        window.pageYOffset || document.documentElement.scrollTop;
      if (scrollTop > lastScrollTop && scrollTop > 55) {
        // If we scroll down more than 50px, we hide the div.
        setVisible(false);
      } else {
        // Otherwise, we display the div again.
        setVisible(true);
      }
      setLastScrollTop(scrollTop);
    };
    window.addEventListener("scroll", handleScroll);
    return () => {
      window.removeEventListener("scroll", handleScroll);
    };
  }, [lastScrollTop]);

  useEffect(() => {
    const handleResize = () => {
      if (window.innerWidth <= 768) {
        setOpenMenu(false);
      } else {
        setOpenMenu(true);
      }
    };

    window.addEventListener("resize", handleResize);
    handleResize(); // Check initial window size

    return () => {
      window.removeEventListener("resize", handleResize);
    };
  }, []);

  const links = [
    { to: "/status", label: "Status" },
    { to: "/daily", label: "Daily" },
    { to: "/compare", label: "Compare" },
    { to: "/search", label: "Search" },
    { to: "/micro", label: "Micro" },
    { to: "/macro", label: "Macro" },
    { to: "/pr", label: "PR" },
  ]

  return (
    <div className={`header flex ${visible ? "visible" : "hidden"}`}>
      <div className="header__container flex">
        <div className="logo__container">
          <a href="/home" className="flex">
            <img src={logo} alt="logo" />
            <span className="hiddenResponsiveMobile">arewefastyet</span>
          </a>
        </div>
        {openMenu ? (
          <>
            <nav>
              <ul className="header__nav flex">
                {links.map((link, index) => (
                  <div key={index}>
                    <li className={index === links.length - 1 ? "" : "liborder--bottom"}>
                      <NavLink
                        className="header__link"
                        activeclassname="active"
                        onClick={handleMenuClick}
                        to={link.to}
                      >
                        {link.label}
                      </NavLink>
                    </li>
                  </div>
                ))}
              </ul>
              <button className="buttonMode" onClick={handleButtonClick}>
                {isColorChanged ? (
                  <i className="fa-solid fa-moon fa-2x " style={{ "color": "#e77002" }}></i>
                ) : (
                  <i className="far fa-sun fa-xl" style={{ "color": "#ffffffdc" }}></i>
                )}
              </button>
            </nav>

            <figure
              className="hamburgerMenu"
              onClick={() => setOpenMenu(false)}
            >
              <figure className="lineMenuTop"></figure>
              <figure className="lineMenuBottom"></figure>
            </figure>
          </>
        ) : (
          <figure className="hamburgerMenu" onClick={() => setOpenMenu(true)}>
            <figure className="lineMenuTop"></figure>
            <figure className="lineMenuBottom"></figure>
          </figure>
        )}
      </div>
    </div>
  );
};

export default Header;
