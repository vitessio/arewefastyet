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

import '../Macro/macro.css'
import "swiper/css";
import "swiper/css/pagination";

import { Mousewheel, Pagination, Keyboard } from "swiper";
import Macrobench from '../../components/Macrobench/Macrobench';
import MacrobenchMobile from '../../components/MacrobenchMobile/MacrobenchMobile';
import { errorApi, openDropDownValue, closeDropDownValue } from '../../utils/utils';

const Macro = () => {
    const urlParams = new URLSearchParams(window.location.search);
    // The following code sets up state variables `gitRefLeft` and `gitRefRight` using the `useState` hook.
    // The values of these variables are based on the query parameters extracted from the URL.

    // If the 'ltag' query parameter is null or undefined, set the initial value of `gitRefLeft` to 'Left',
    // otherwise, use the value of the 'ltag' query parameter.
    const [gitRefLeft, setGitRefLeft] = useState(urlParams.get('ltag') == null ? 'Left' : urlParams.get('ltag'));
    const [gitRefRight, setGitRefRight] = useState(urlParams.get('rtag') == null ? 'Right' : urlParams.get('rtag'));
    const [openDropDownLeft, setOpenDropDownLeft] = useState(openDropDownValue);
    const [openDropDownRight, setOpenDropDownRight] = useState(openDropDownValue);
    const [dataRefs, setDataRefs] = useState([]);
    const [isFirstCallFinished,setIsFirstCallFinished] = useState(false)
    const [dataMacrobench, setDataMacrobench] = useState([]); 
    const [commitHashLeft, setCommitHashLeft] = useState('')
    const [commitHashRight, setCommitHashRight] = useState('')
    const [error, setError] = useState(null);
    const [isLoading, setIsLoading] = useState(true);
    const [currentSlideIndex, setCurrentSlideIndex] = useState(urlParams.get('ptag') == null ? '0' : urlParams.get('ptag'));
    const [currentSlideIndexMobile, setCurrentSlideIndexMobile] = useState(urlParams.get('ptagM') == null ? '0' : urlParams.get('ptagM'))
    
    // updateCommitHash: This function updates the value of CommitHash based on the provided Git reference and JSON data.
    const updateCommitHash = (gitRef, setCommitHash, jsonDataRefs) => {
        const obj = jsonDataRefs.find(item => item.Name === gitRef);
            setCommitHash(obj ? obj.CommitHash : null);
    }
    
    useEffect(() => {
        const fetchData = async () => {
          try {
            const responseRefs = await fetch(`${import.meta.env.VITE_API_URL}vitess/refs`);
      
            const jsonDataRefs = await responseRefs.json();
            
            setDataRefs(jsonDataRefs);
            setIsLoading(false)

            updateCommitHash(gitRefLeft, setCommitHashLeft, jsonDataRefs)
            updateCommitHash(gitRefRight, setCommitHashRight, jsonDataRefs)
            
            setIsFirstCallFinished(true)

          } catch (error) {
            console.log('Error while retrieving data from the API', error);
            setError(errorApi);
            setIsLoading(false);
          }
        };
      
        fetchData();
      }, []);

    useEffect(() => {
        if (isFirstCallFinished) {
            const fetchData = async () => {
                try {
                    const responseMacrobench = await fetch(`${import.meta.env.VITE_API_URL}macrobench/compare?rtag=${commitHashRight}&ltag=${commitHashLeft}`)

                    const jsonDataMacrobench = await responseMacrobench.json();
                    setDataMacrobench(jsonDataMacrobench)
                } catch (error) {
                    console.log('Error while retrieving data from the API', error);
                    setError(errorApi);
                }
            };
            fetchData();
        }  
    }, [commitHashLeft, commitHashRight])
       

    // OPEN DROP DOWN

    const openDropDown = (openDropDown, setOpenDropDown) =>{
        if (openDropDown === closeDropDownValue) {
            setOpenDropDown(openDropDownValue);
          } else {
            setOpenDropDown(closeDropDownValue);
          }
    }

    // CHANGE VALUE DROPDOWN

    const valueDropDown = (ref, setDropDown, setCommitHash, setOpenDropDown, setChangeUrl) => {
        setDropDown(ref.Name)
        setCommitHash(ref.CommitHash)
        setOpenDropDown(closeDropDownValue);
    }

    // Changing the URL relative to the reference of a selected benchmark.
    // Storing the carousel position as a URL parameter.
    const navigate = useNavigate();
    
    useEffect(() => {
        navigate(`?ltag=${gitRefLeft}&rtag=${gitRefRight}&ptag=${currentSlideIndex}&ptagM=${currentSlideIndexMobile}`)
    }, [gitRefLeft, gitRefRight, currentSlideIndex, currentSlideIndexMobile]) 
    
    

    const handleSlideChange = (swiper) => {
        setCurrentSlideIndex(swiper.realIndex);
    };
        
    return (
        <div className='macro'>
            <div className='macro__top justify--content'>
                <div className='macro__top__text'>
                    <h2>Compare Macrobenchmarks</h2>
                    <span>
                            Lorem ipsum dolor sit amet, consectetur adipiscing elit. Integer a augue mi.
                            Etiam sed imperdiet ligula, vel elementum velit.
                            Phasellus sodales felis eu condimentum convallis.
                            Suspendisse sodales malesuada iaculis. Mauris molestie placerat ex non malesuada.
                            Curabitur eget sagittis eros. Aliquam aliquam sem non tincidunt volutpat. 
                    </span>
                    
                </div>
                <figure className='macroStats'></figure>
            </div>
            <figure className='line'></figure>
            <div className='macro__bottom'>
                <h3>Comparing result for <a target="_blank" href={commitHashLeft ? `https://github.com/vitessio/vitess/commit/${commitHashLeft}` : undefined}>{gitRefLeft}</a> and <a target="_blank" href={commitHashRight ? `https://github.com/vitessio/vitess/commit/${commitHashRight}` : undefined}>{gitRefRight}</a> </h3>
                <div className='macro__bottom__DropDownContainer'>
                    <figure className='macro__bottom__DropDownLeft flex--column' style={{ maxHeight: `${openDropDownLeft}px` }}>
                        <span className='DropDown__Base'  onClick={() => openDropDown(openDropDownLeft, setOpenDropDownLeft)}>{gitRefLeft} <i className="fa-solid fa-circle-arrow-down"></i></span>
                        {dataRefs.map((ref, index) => {
                            return (
                                <React.Fragment key={index}>
                                    <figure className='dropDown--line'></figure>
                                    <span className='dropDown__ref' onClick={() => {valueDropDown(ref, setGitRefLeft, setCommitHashLeft, setOpenDropDownLeft)}}>{ref.Name}</span>
                                </React.Fragment>
                            )
                        })}
                    </figure>
                    <figure className='macro__bottom__DropDownRight flex--column' style={{ maxHeight: `${openDropDownRight}px` }}>
                        <span className='DropDown__Base'  onClick={() => openDropDown(openDropDownRight, setOpenDropDownRight)}>{gitRefRight} <i className="fa-solid fa-circle-arrow-down"></i></span>
                        {dataRefs.map((ref, index) => {
                            return (
                                <React.Fragment key={index}>
                                    <figure className='dropDown--line'></figure>
                                    <span className='dropDown__ref' onClick={() => valueDropDown(ref, setGitRefRight, setCommitHashRight, setOpenDropDownRight)}>{ref.Name}</span>
                                </React.Fragment>
                            )
                        })}
                    </figure>
                </div>
                {error ? (
                    <div className='macrobench__apiError'>{error}</div> 
                ) : (
                    isLoading ? (
                        <div className='loadingSpinner'>
                            <RingLoader loading={isLoading} color='#E77002' size={300}/>
                            </div>
                        ): ( 
                            <div className='macrobench__Container flex'>
                                <div className='macrobench__Sidebar flex--column'>
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
                                <div className='carousel__container'>
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

export default Macro;