import React from 'react';

import logo from '../../assets/logo.png'
import homeLogo from '../../assets/homeLogoLarge.png'

import './home.css'

const Home = () => {
    return (
        <div className='home'>
            <article className='home__top justify--content'>
                <div className='home__top__text'>
                    <h1>Are We Fast Yet</h1>
                    <span>
                        Lorem ipsum dolor sit amet, consectetur adipiscing elit. Integer a augue mi.
                        Etiam sed imperdiet ligula, vel elementum velit.
                        Phasellus sodales felis eu condimentum convallis.
                        Suspendisse sodales malesuada iaculis. Mauris molestie placerat ex non malesuada.
                        Curabitur eget sagittis eros. Aliquam aliquam sem non tincidunt volutpat.
                    </span>
                </div>
                <img src={homeLogo} className='home__top__logo'/>
            
            </article>
            <figure className='line'></figure>
            <article className='home__bottom'>
                <span > 
                    Lorem ipsum dolor sit amet, consectetur adipiscing elit.
                    Integer a augue mi.
                    Etiam sed imperdiet ligula, vel elementum velit.
                    Phasellus sodales felis eu condimentum convallis.
                    Suspendisse sodales malesuada iaculis. Mauris molestie placerat ex non malesuada.
                    Curabitur eget sagittis eros. Aliquam aliquam sem non tincidunt volutpat.
                    Ut sodales ut justo a rutrum. Proin ac nunc sem. Aenean varius vestibulum tortor, eget lacinia massa malesuada ut.
                    Vivamus dolor justo, rhoncus eget risus eu, lobortis convallis justo.
                    Nunc imperdiet imperdiet ante vel pharetra.
                    Fusce ut arcu sollicitudin, posuere odio eget, lobortis leo.
                    Nulla eget libero nisi.
                </span>
            </article>
        </div>
    );
};

export default Home;