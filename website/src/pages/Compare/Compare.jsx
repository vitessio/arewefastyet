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

import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import RingLoader from "react-spinners/RingLoader";
import { Swiper, SwiperSlide } from "swiper/react";

import '../Compare/compare.css'
import "swiper/css";
import "swiper/css/pagination";

import { Mousewheel, Pagination, Keyboard } from "swiper";
import { errorApi, updateCommitHash} from '../../utils/Utils';
import Macrobench from '../../components/Macrobench/Macrobench';
import MacrobenchMobile from '../../components/MacrobenchMobile/MacrobenchMobile';

const Compare = () => {

    const urlParams = new URLSearchParams(window.location.search);
    // The following code sets up state variables `gitRefLeft` and `gitRefRight` using the `useState` hook.
    // The values of these variables are based on the query parameters extracted from the URL.

    // If the 'ltag' query parameter is null or undefined, set the initial value of `gitRefLeft` to 'Left',
    // otherwise, use the value of the 'ltag' query parameter.
    const [gitRefLeft, setGitRefLeft] = useState(urlParams.get('ltag') == null ? 'Left' : urlParams.get('ltag'))
    const [gitRefRight, setGitRefRight] = useState(urlParams.get('rtag') == null ? 'Right' : urlParams.get('rtag'))
    const [error, setError] = useState(null);
    const [isLoading, setIsLoading] = useState(true);
    const [dataMacrobench, setDataMacrobench] = useState([]); 
    const [currentSlideIndex, setCurrentSlideIndex] = useState(urlParams.get('ptag') == null ? '0' : urlParams.get('ptag'));
    const [currentSlideIndexMobile, setCurrentSlideIndexMobile] = useState(urlParams.get('ptagM') == null ? '0' : urlParams.get('ptagM'))

    useEffect(() => {
            const fetchData = async () => {
                try {
                    const responseMacrobench = await fetch(`${import.meta.env.VITE_API_URL}macrobench/compare?rtag=${gitRefRight}&ltag=${gitRefLeft}`)

                    const jsonDataMacrobench = await responseMacrobench.json();
                    console.log(jsonDataMacrobench)
                    setDataMacrobench(jsonDataMacrobench)
                    setIsLoading(false)
                } catch (error) {
                    console.log('Error while retrieving data from the API', error);
                    setError(errorApi);
                    setIsLoading(false)
                }
            };
            fetchData();
        
    }, [gitRefLeft, gitRefRight])

    // Changing the URL relative to the reference of a selected benchmark.
    // Storing the carousel position as a URL parameter.
    const navigate = useNavigate();
    
    useEffect(() => {
        navigate(`?ltag=${gitRefLeft}&rtag=${gitRefRight}`)
    }, [gitRefLeft, gitRefRight]) 
    
    const handleInputChangeLeft = (e) => {
        setGitRefLeft(e.target.value);
      };
    

      const handleInputChangeRight = (e) => {
        setGitRefRight(e.target.value);
      };

      const handleSlideChange = (swiper) => {
        setCurrentSlideIndex(swiper.realIndex);
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

                {error ? (
                    <div className='macrobench__apiError'>{error}</div> 
                ) : (
                    isLoading ? (
                        <div className='loadingSpinner'>
                            <RingLoader loading={isLoading} color='#E77002' size={300}/>
                            </div>
                        ): ( 
                            <div className='compare__macrobench__Container flex'>
                                <div className='compare__macrobench__Sidebar flex--column'>
                                    <span >QPS Total</span>
                                    <figure className='macrobench__Sidebar__line'></figure>
                                    <span>QPS Reads</span>
                                    <figure className='macrobench__Sidebar__line'></figure>
                                    <span>QPS Writes</span>
                                    <figure className='macrobench__Sidebar__line'></figure>
                                    <span>QPS Other</span>
                                    <figure className='macrobench__Sidebar__line'></figure>
                                    <span>TPS</span>
                                    <figure className='macrobench__Sidebar__line'></figure>
                                    <span>Latency</span>
                                    <figure className='macrobench__Sidebar__line'></figure>
                                    <span>Errors</span>
                                    <figure className='macrobench__Sidebar__line'></figure>
                                    <span>Reconnects</span>
                                    <figure className='macrobench__Sidebar__line'></figure>
                                    <span>Time</span>
                                    <figure className='macrobench__Sidebar__line'></figure>
                                    <span>Threads</span>
                                    <figure className='macrobench__Sidebar__line'></figure>
                                    <span>Total CPU time</span>
                                    <figure className='macrobench__Sidebar__line'></figure>
                                    <span>CPU time vtgate</span>
                                    <figure className='macrobench__Sidebar__line'></figure>
                                    <span>CPU time vttablet</span>
                                    <figure className='macrobench__Sidebar__line'></figure>
                                    <span>Total Allocs bytes</span>
                                    <figure className='macrobench__Sidebar__line'></figure>
                                    <span>Allocs bytes vtgate</span>
                                    <figure className='macrobench__Sidebar__line'></figure>
                                    <span>Allocs bytes vttablet</span>
                                </div>
                                <div className='compare__carousel__container'>
                                    <Swiper
                                        direction={"vertical"}
                                        slidesPerView={1}
                                        spaceBetween={30}
                                        mousewheel={true}
                                        keyboard={{
                                            enabled: true,
                                        }}
                                        pagination={{
                                        clickable: true,
                                        }}
                                        modules={[Mousewheel, Pagination, Keyboard]}
                                        onSlideChange={handleSlideChange}
                                        initialSlide={currentSlideIndex}
                                        className="mySwiper"
                                        >
                                        {dataMacrobench.map((macro, index) => {
                                            return (
                                                <SwiperSlide key={index}>
                                                    <Macrobench data={macro} gitRefLeft={gitRefLeft} gitRefRight={gitRefRight} swiperSlide={SwiperSlide}/>
                                                    <MacrobenchMobile 
                                                    data={macro} 
                                                    gitRefLeft={gitRefLeft} 
                                                    gitRefRight={gitRefRight} 
                                                    swiperSlide={SwiperSlide} 
                                                    handleSlideChange={handleSlideChange} 
                                                    setCurrentSlideIndexMobile={setCurrentSlideIndexMobile}
                                                    currentSlideIndexMobile={currentSlideIndexMobile}
                                                    />
                                                </SwiperSlide>
                                            )
                                        })}
                                        
                                    </Swiper>
                                </div>                
                            </div>
                    ))}
                
            </div>
        </div>
    );
};

export default Compare;