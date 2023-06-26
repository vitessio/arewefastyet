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

import React, { useState } from 'react';

import '../Compare/compare.css'

const Compare = () => {

    const [gitRefLeft, setGitRefLeft] = useState('')
    const [gitRefRight, setGitRefRight] = useState('')
    
    const handleInputChangeLeft = (e) => {
        setGitRefLeft(e.target.value);
      };
    

      const handleInputChangeRight = (e) => {
        setGitRefRight(e.target.value);
      };
    return (
        <div className='compare'>
            <div className='compare__top justify--content'>
                <div className='compare__top__text'>
                    <h2>Compare</h2>
                    <span>
                            Lorem ipsum dolor sit amet, consectetur adipiscing elit. Integer a augue mi.
                            Etiam sed imperdiet ligula, vel elementum velit.
                            Phasellus sodales felis eu condimentum convallis.
                            Suspendisse sodales malesuada iaculis. Mauris molestie placerat ex non malesuada.
                            Curabitur eget sagittis eros. Aliquam aliquam sem non tincidunt volutpat. 
                    </span>
                    
                </div>
                <figure className='compareStats'></figure>
            </div>
            <figure className='line'></figure>
            <div className='compare__bottom'>
                <div className='justify--content form__container'>
                    <h3>Comparing <a href={`https://github.com/vitessio/vitess/commit/${gitRefLeft}`}>{gitRefLeft.slice(0, 8)}</a> with <a href={`https://github.com/vitessio/vitess/commit/${gitRefRight}`}>{gitRefRight.slice(0, 8)}</a></h3>
                    <form className='justify--content'>
                        <input
                        type="text"
                        value={gitRefLeft}
                        onChange={handleInputChangeLeft}
                        placeholder="Left commit SHA"
                        className='form__inputLeft'></input>
                        <input
                        type="text"
                        value={gitRefRight}
                        onChange={handleInputChangeRight}
                        placeholder="Left commit SHA"
                        className='form__inputRight'
                        >
                        </input>
                        <button>Compare</button>
                    </form>
                </div>
                
            </div>
        </div>
    );
};

export default Compare;