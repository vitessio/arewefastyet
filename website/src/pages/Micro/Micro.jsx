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

import React, {useState, useEffect} from 'react';
import {useNavigate} from 'react-router-dom';

import '../Micro/micro.css'

import { errorApi, openDropDownValue, closeDropDownValue } from '../../utils/utils';
import Microbench from '../../components/Microbench/Microbench';

const Micro = () => {
    const urlParams = new URLSearchParams(window.location.search);
    const [gitRefLeft, setGitRefLeft] = useState(urlParams.get('ltag') == null ? 'Left' : urlParams.get('ltag'));
    const [gitRefRight, setGitRefRight] = useState(urlParams.get('rtag') == null ? 'Right' : urlParams.get('rtag'));
    const [openDropDownLeft, setOpenDropDownLeft] = useState(closeDropDownValue);
    const [openDropDownRight, setOpenDropDownRight] = useState(closeDropDownValue);
    const [dataRefs, setDataRefs] = useState([]);
    const [dataMicrobench, setDataMicrobench] = useState([]); 
    const [isFirstCallFinished,setIsFirstCallFinished] = useState(false)
    const [commitHashLeft, setCommitHashLeft] = useState('')
    const [commitHashRight, setCommitHashRight] = useState('')
    const [error, setError] = useState(null);
    const [isLoading, setIsLoading] = useState(true);

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
      console.log(commitHashLeft)
      useEffect(() => {
        if (isFirstCallFinished) {
            const fetchData = async () => {
                try {
                    const responseMicrobench = await fetch(`${import.meta.env.VITE_API_URL}microbench/compare?rtag=83eef5d198ae1d9c5245270c8ca39961ca847b13&ltag=60dc8b997fccb0060ed2b29f811dc0dea38e9123`)

                    const jsonDataMicrobench = await responseMicrobench.json();
                    setDataMicrobench(jsonDataMicrobench)
                    console.log(jsonDataMicrobench)
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
        navigate(`?ltag=${gitRefLeft}&rtag=${gitRefRight}`)
    }, [gitRefLeft, gitRefRight]) 


    return (
        <div className='micro'>
            <div className='micro__top justify--content'>
                <div className='micro__top__text'>
                    <h2>MicroBenchmark</h2>
                    <span>
                            Lorem ipsum dolor sit amet, consectetur adipiscing elit. Integer a augue mi.
                            Etiam sed imperdiet ligula, vel elementum velit.
                            Phasellus sodales felis eu condimentum convallis.
                            Suspendisse sodales malesuada iaculis. Mauris molestie placerat ex non malesuada.
                            Curabitur eget sagittis eros. Aliquam aliquam sem non tincidunt volutpat. 
                    </span>
                </div>
                <figure className='microStats'></figure>
            </div>
            <figure className='line'></figure>
            <div className='micro__bottom'>
                <div className='micro__bottom__title justify--content'>
                    <h3>Comparing result for</h3>
                    <figure className='micro__bottom__DropDownLeft flex--column' style={{ maxHeight: `${openDropDownLeft}px` }}>
                        <span className='DropDown__BaseSpan'  onClick={() => openDropDown(openDropDownLeft, setOpenDropDownLeft)}>{gitRefLeft} <i className="fa-solid fa-circle-arrow-down"></i></span>
                        {dataRefs.map((ref, index) => {
                            return (
                                <React.Fragment key={index}>
                                    <figure className='dropDown--line'></figure>
                                    <span className='dropDown__ref' onClick={() => {valueDropDown(ref, setGitRefLeft, setCommitHashLeft, setOpenDropDownLeft)}}>{ref.Name}</span>
                                </React.Fragment>
                            )
                        })}
                    </figure>
                    <h3>and</h3>
                    <figure className='micro__bottom__DropDownRight flex--column' style={{ maxHeight: `${openDropDownRight}px` }}>
                        <span className='DropDown__BaseSpan'  onClick={() => openDropDown(openDropDownRight, setOpenDropDownRight)}>{gitRefRight} <i className="fa-solid fa-circle-arrow-down"></i></span>
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
                <div className='micro__container'>
                    <div className='micro__thead space--between'>
                        <span className='width--12em'>Package</span>
                        <span className='width--14em'>Benchmark Name</span>
                        <span className='width--18em'>Number of Iterations</span>
                        <span className='width--18em'>Time/op</span>
                        <span className='width--6em'>More</span>
                    </div>
                    <figure className='micro__thead__line'></figure>
                    <div className='space--between data__top'>
                        <div className='width--12em'></div>
                        <div className='width--14em'></div>
                        <div className='width--18em space--between'>
                            <span className='width--100'>{gitRefLeft}</span> 
                            <span className='width--100'>{gitRefRight}</span>
                            <span className='width--100'>Diff %</span>
                        </div>
                        <div className='width--18em space--between'>
                            <span className='width--100'>{gitRefLeft}</span> 
                            <span className='width--100'>{gitRefRight}</span>
                            <span className='width--100'>Diff %</span>
                        </div>
                        <div className='width--6em'></div>
                    </div>
                    {dataMicrobench.map((micro, index) => {
                        const isEvenIndex = index % 2 === 0;
                        const backgroundGrey = isEvenIndex ? 'grey--background' : '';
                        return(
                            <Microbench data={micro} key={index} className={backgroundGrey} gitRefLeft={gitRefLeft} gitRefRight={gitRefRight}/>
                        )
                    })}
                </div>
            </div>
        </div>
    );
};

export default Micro;