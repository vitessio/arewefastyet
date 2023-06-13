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
import RingLoader from "react-spinners/RingLoader";

import '../Macro/macro.css'
import Macrobench from '../../components/Macrobench/Macrobench';

const Macro = () => {
    const [dropDownLeft, setDropDownLeft] = useState('Left');
    const [dropDownRight, setDropDownRight] = useState('Right');
    const [openDropDownLeft, setOpenDropDownLeft] = useState(58);
    const [openDropDownRight, setOpenDropDownRight] = useState(58);
    const [dataRefs, setDataRefs] = useState([]);
    const [dataMacrobench, setDataMacrobench] = useState([]);
    const [commitHashLeft, setCommitHashLeft] = useState('')
    const [commitHashRight, setCommitHashRight] = useState('')
    const [error, setError] = useState(null);
    const [isLoading, setIsLoading] = useState(true);
    
    

    useEffect(() => {
        const fetchData = async () => {
          try {
            const responseRefs = await fetch(`${import.meta.env.VITE_API_URL}vitess/refs`);
      
            const jsonDataRefs = await responseRefs.json();
            
            setDataRefs(jsonDataRefs);
            console.log(jsonDataRefs)
            setIsLoading(false)
          } catch (error) {
            console.log('Error while retrieving data from the API', error);
            setError('An error occurred while retrieving data from the API. Please try again.');
            setIsLoading(false);
          }
        };
      
        fetchData();
      }, []);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const responseMacrobench = await fetch(`${import.meta.env.VITE_API_URL}macrobench/compare?rtag=${commitHashRight}&ltag=${commitHashLeft}`)

                const jsonDataMacrobench = await responseMacrobench.json();
                console.log(jsonDataMacrobench)
                setDataMacrobench(jsonDataMacrobench)
            } catch (error) {
                console.log('Error while retrieving data from the API', error);
                setError('An error occurred while retrieving data from the API. Please try again.');
            }
        };
        fetchData();
    }, [commitHashLeft, commitHashRight])
       
        

    


    // OPEN DROP DOWN

    const openDropDown = (openDropDown, setOpenDropDown) =>{
        if (openDropDown === 58) {
            setOpenDropDown(1000);
          } else {
            setOpenDropDown(58);
          }
    }

    // CHANGE VALUE DROPDOWN

    const valueDropDown = (ref, setDropDown, setCommitHash, setOpenDropDown) => {
        setDropDown(ref.Name)
        setCommitHash(ref.CommitHash)
        setOpenDropDown(58);
    }

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
                <h3>Comparing result for <a target="_blank" href={commitHashLeft ? `https://github.com/vitessio/vitess/commit/${commitHashLeft}` : undefined}>{dropDownLeft}</a> and <a target="_blank" href={commitHashRight ? `https://github.com/vitessio/vitess/commit/${commitHashRight}` : undefined}>{dropDownRight}</a> </h3>
                <div className='macro__bottom__DropDownContainer'>
                    <figure className='macro__bottom__DropDownLeft flex--column' style={{ maxHeight: `${openDropDownLeft}px` }}>
                        <span className='DropDown__Base'  onClick={() => openDropDown(openDropDownLeft, setOpenDropDownLeft)}>{dropDownLeft} <i className="fa-solid fa-circle-arrow-down"></i></span>
                        {dataRefs.map((ref, index) => {
                            return (
                                <React.Fragment key={index}>
                                    <figure className='dropDown--line'></figure>
                                    <span className='dropDown__ref' onClick={() => valueDropDown(ref, setDropDownLeft, setCommitHashLeft, setOpenDropDownLeft)}>{ref.Name}</span>
                                </React.Fragment>
                            )
                        })}
                    </figure>
                    <figure className='macro__bottom__DropDownRight flex--column' style={{ maxHeight: `${openDropDownRight}px` }}>
                        <span className='DropDown__Base'  onClick={() => openDropDown(openDropDownRight, setOpenDropDownRight)}>{dropDownRight} <i className="fa-solid fa-circle-arrow-down"></i></span>
                        {dataRefs.map((ref, index) => {
                            return (
                                <React.Fragment key={index}>
                                    <figure className='dropDown--line'></figure>
                                    <span className='dropDown__ref' onClick={() => valueDropDown(ref, setDropDownRight, setCommitHashRight, setOpenDropDownRight)}>{ref.Name}</span>
                                </React.Fragment>
                            )
                        })}
                    </figure>
                </div>

                <div className='macrobench__Container flex'>
                    <div className='macrobench__Sidebar flex--column'>
                        <span>QPS Total</span>
                        <span>QPS Reads</span>
                        <span>QPS Writes</span>
                        <span>QPS Other</span>
                        <span>TPS</span>
                        <span>Latency</span>
                        <span>Errors</span>
                        <span>Reconnects</span>
                        <span>Time</span>
                        <span>Threads</span>
                        <span>Total CPU time</span>
                        <span>CPU time vtgate</span>
                        <span>CPU time vttablet</span>
                        <span>Total Allocs bytes</span>
                        <span>Allocs bytes vtgate</span>
                        <span>Allocs bytes vttablet</span>
                    </div>
                    <div className='macrobench__Compare'>
                        {dataMacrobench.map((macro, index) => {
                            return (
                                <Macrobench key={index} data={macro} dropDownLeft={dropDownLeft} dropDownRight={dropDownRight}/>
                            )
                        })}
                    </div>
                </div>
            </div>
        </div>
    );
};

export default Macro;