import React from 'react';

import logo from '../../assets/logo.png'

import './home.css'

const Home = () => {
    return (
        <div className='home'>
            <article className='home__top'>
                <div className='home__top__text'>
                    <h1>Are We Fast Yet</h1>
                    <span>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Integer a augue mi.
                         Etiam sed imperdiet ligula, vel elementum velit.
                          Phasellus sodales felis eu condimentum convallis.
                           Suspendisse sodales malesuada iaculis. Mauris molestie placerat ex non malesuada.
                            Curabitur eget sagittis eros. Aliquam aliquam sem non tincidunt volutpat.
                    </span>
                </div>
                <div className='logo__container'>
                    <img src={logo}/>
                </div>

            </article>
        </div>
    );
};

export default Home;