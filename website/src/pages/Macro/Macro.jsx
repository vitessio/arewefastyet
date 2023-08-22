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
import { SwiperSlide } from "swiper/react";

import '../Macro/Macro.css'
// import "swiper/css";


import Macrobench from '../../components/MacroComponents/Macrobench/Macrobench';
import MacrobenchMobile from '../../components/MacroComponents/MacrobenchMobile/MacrobenchMobile';
import { errorApi, closeDropDownValue, updateCommitHash, openDropDown, valueDropDown } from '../../utils/Utils';

const Macro = () => {
    const urlParams = new URLSearchParams(window.location.search);
    // The following code sets up state variables `gitRefLeft` and `gitRefRight` using the `useState` hook.
    // The values of these variables are based on the query parameters extracted from the URL.

    // If the 'ltag' query parameter is null or undefined, set the initial value of `gitRefLeft` to 'Left',
    // otherwise, use the value of the 'ltag' query parameter.
    const [gitRefLeft, setGitRefLeft] = useState(urlParams.get('ltag') == null ? 'Left' : urlParams.get('ltag'));
    const [gitRefRight, setGitRefRight] = useState(urlParams.get('rtag') == null ? 'Right' : urlParams.get('rtag'));
    const [openDropDownLeft, setOpenDropDownLeft] = useState(closeDropDownValue);
    const [openDropDownRight, setOpenDropDownRight] = useState(closeDropDownValue);
    const [dataRefs, setDataRefs] = useState([]);
    const [isFirstCallFinished,setIsFirstCallFinished] = useState(false)
    const [dataMacrobench, setDataMacrobench] = useState([]); 
    const [commitHashLeft, setCommitHashLeft] = useState('')
    const [commitHashRight, setCommitHashRight] = useState('')
    const [error, setError] = useState(null);
    const [isLoading, setIsLoading] = useState(true);
    const [currentSlideIndexMobile, setCurrentSlideIndexMobile] = useState(urlParams.get('ptagM') == null ? '0' : urlParams.get('ptagM'))
    const [textLoading, setTextLoading] = useState(true)
    
    useEffect(() => {
        const fetchData = async () => {
          try {
            const responseRefs = await fetch(`${import.meta.env.VITE_API_URL}vitess/refs`);
      
            const jsonDataRefs = await responseRefs.json();
            
            setDataRefs(jsonDataRefs);
            

            updateCommitHash(gitRefLeft, setCommitHashLeft, jsonDataRefs)
            updateCommitHash(gitRefRight, setCommitHashRight, jsonDataRefs)
            
            setIsFirstCallFinished(true)

          } catch (error) {
            console.log('Error while retrieving data from the API', error);
            setError(errorApi);
            
          }
        };
      
        fetchData();
      }, []);

    useEffect(() => {
        if (isFirstCallFinished) {
            setTextLoading(true)
            const fetchData = async () => {
                try {
                    const responseMacrobench = await fetch(`${import.meta.env.VITE_API_URL}macrobench/compare?rtag=${commitHashRight}&ltag=${commitHashLeft}`)

                    const jsonDataMacrobench = await responseMacrobench.json();
                    setDataMacrobench(jsonDataMacrobench)
                    setIsLoading(false)
                    setTextLoading(false)
                } catch (error) {
                    console.log('Error while retrieving data from the API', error);
                    setError(errorApi);
                    setIsLoading(false);
                    setTextLoading(false)
                }
            };
            fetchData();
        }  
    }, [commitHashLeft, commitHashRight])


    // Changing the URL relative to the reference of a selected benchmark.
    // Storing the carousel position as a URL parameter.
    const navigate = useNavigate();
    
    useEffect(() => {
        navigate(`?ltag=${gitRefLeft}&rtag=${gitRefRight}&ptagM=${currentSlideIndexMobile}`)
    }, [gitRefLeft, gitRefRight, currentSlideIndexMobile]) 
    

    const handleSlideChange = (swiper) => {
        setCurrentSlideIndexMobile(swiper.realIndex);
    };
        
    return (
        <div className='macro'>
            <div className='macro__top justify--content'>
                <div className='macro__top__text'>
                    <h2 className='header--title'>Compare Macrobenchmarks</h2>
                    <div className='macro__bottom__DropDownContainer flex'>
                    <figure className='macro__bottom__DropDownLeft dropDown flex--column' style={{ maxHeight: `${openDropDownLeft}px` }}>
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
                    <figure className='macro__bottom__DropDownRight dropDown flex--column' style={{ maxHeight: `${openDropDownRight}px` }}>
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
                    
                </div>
                <figure className='macroStats'></figure>
            </div>
            <figure className='line'></figure>
            <div className='macro__bottom'>
               
                
                {error ? (
                    <div className='apiError'>{error}</div> 
                ) : (
                    isLoading ? (
                        <div className='loadingSpinner'>
                            <RingLoader loading={isLoading} color='#E77002' size={300}/>
                            </div>
                        ): ( 
                            <div className='macrobench__Container flex'>
                                
                                <div className='carousel__container'>
                                    
                                        {dataMacrobench.map((macro, index) => {
                                            return (
                                                <div key={index}>
                                                    <Macrobench
                                                        data={macro} 
                                                        gitRefLeft={gitRefLeft} 
                                                        gitRefRight={gitRefRight} 
                                                        swiperSlide={SwiperSlide} 
                                                        commitHashLeft={commitHashLeft}
                                                        commitHashRight={commitHashRight}
                                                        textLoading={textLoading}
                                                    />
                                                    <MacrobenchMobile 
                                                    data={macro} 
                                                    gitRefLeft={gitRefLeft} 
                                                    gitRefRight={gitRefRight} 
                                                    swiperSlide={SwiperSlide} 
                                                    handleSlideChange={handleSlideChange} 
                                                    setCurrentSlideIndexMobile={setCurrentSlideIndexMobile}
                                                    currentSlideIndexMobile={currentSlideIndexMobile}
                                                    commitHashLeft={commitHashLeft}
                                                    commitHashRight={commitHashRight}
                                                    textLoading={textLoading}
                                                    />
                                                </div>
                                            )
                                        })}
                                        
                                    
                                </div>                
                            </div>
                    ))}
            </div>
                
        </div>
    );
};

export default Macro;