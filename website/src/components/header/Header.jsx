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

import React from 'react';
import {NavLink} from 'react-router-dom';
import './header.css'

//import image

import logo from '../../assets/logo.png'

const Header = () => {
    return (
        <div className='header flex'>
            <div className='logo__container justify--content'>
                <img src={logo} alt='logo'/>
                <span>Benchmark</span>
            </div>
            
            <nav>
                <ul className='header__nav flex'>
                    {[
                        { to: '/home', label: 'Home' },
                        { to: '/status', label: 'Status' },
                        { to: '/status', label: 'CRON' },
                        { to: '/status', label: 'Compare' },
                        { to: '/status', label: 'Search' },
                        { to: '/status', label: 'Micro' },
                        { to: '/status', label: 'Macro' }
                        ].map((link, index) => (
                        <li key={index}>
                            <NavLink className='header__link' activeclassname='active' to={link.to}>
                            {link.label}
                            </NavLink>
                        </li>
                        ))}
                </ul>
            </nav>
            
        </div>
    );
};

export default Header;