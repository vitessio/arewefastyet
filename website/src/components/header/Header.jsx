import React from 'react';
import {NavLink} from 'react-router-dom';
import './header.css'

//import image

import logo from '../../assets/logo.png'

const Header = () => {
    return (
        <div className='header flex'>
            <div className='logo__container justify--content'>
                <img src={logo}/>
                <span>Benchmark</span>
            </div>
            
            <nav>
                <ul className='header__nav flex'>
                    <li><NavLink activeclassname='active' to='/home'>Home</NavLink></li>
                    <li><NavLink to='/status'>Status</NavLink></li>
                    <li><NavLink to='/status'>CRON</NavLink></li>
                    <li><NavLink to='/status'>Compare</NavLink></li>
                    <li><NavLink to='/status'>Search</NavLink></li>
                    <li><NavLink to='/status'>Micro</NavLink></li>
                    <li><NavLink to='/status'>Macro</NavLink></li>
                </ul>
            </nav>
            
        </div>
    );
};

export default Header;