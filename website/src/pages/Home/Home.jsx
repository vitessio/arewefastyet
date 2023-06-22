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

import logo from '../../assets/logo.png'
import homeLogo from '../../assets/homeLogoLarge.png'

import './home.css'

const Home = () => {
    return (
        <div className='home'>
            <article className='home__top justify--content'>
                <div className='home__top__text'>
                    <h1>AreWeFastYet</h1>
                    <span>
                        Lorem ipsum dolor sit amet, consectetur adipiscing elit. Integer a augue mi.
                        Etiam sed imperdiet ligula, vel elementum velit.
                        Phasellus sodales felis eu condimentum convallis.
                        Suspendisse sodales malesuada iaculis. Mauris molestie placerat ex non malesuada.
                        Curabitur eget sagittis eros. Aliquam aliquam sem non tincidunt volutpat.
                    </span>
                </div>
                <img src={homeLogo} alt='logo' className='home__top__logo'/>
            
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