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

import React, {useEffect, useState} from 'react';
import { useNavigate } from 'react-router-dom';
import RingLoader from "react-spinners/RingLoader";
import { Swiper, SwiperSlide } from "swiper/react";

import '../Search/search.css'
import "swiper/css";
import "swiper/css/pagination";

import { errorApi } from '../../utils/Utils';
import { Mousewheel, Pagination, Keyboard } from "swiper";
import SearchMacro from '../../components/SearchMacro/SearchMacro';

const Search = () => {
    const urlParams = new URLSearchParams(window.location.search);
    const [gitRef, setGitRef] = useState(urlParams.get('git_ref') || '');
    const [dataSearch, setDataSearch] = useState([]); 
    const [isFormSubmitted, setIsFormSubmitted] = useState(false);
    const [error, setError] = useState(null);
    const [isLoading, setIsLoading] = useState(true);
    const [currentSlideIndex, setCurrentSlideIndex] = useState(urlParams.get('ptag') == null ? '0' : urlParams.get('ptag'));

    useEffect(() => {
        const fetchData = async () => {
            try {
                const responseSearch = await fetch(`${import.meta.env.VITE_API_URL}search?git_ref=${gitRef}`);
          
                const jsonDataSearch = await responseSearch.json();
                
                setDataSearch(jsonDataSearch);
                console.log(jsonDataSearch)
                setIsLoading(false)
    
              } catch (error) {
                console.log('Error while retrieving data from the API', error);
                setError(errorApi);
                setIsLoading(false);
              }
        }

        fetchData();
    }, [isFormSubmitted])

    const navigate = useNavigate();
    
    useEffect(() => {
        navigate(`?git_ref=${gitRef}`)
    }, [gitRef]) 

    const handleInputChange = (e) => {
        setGitRef(e.target.value);
      };
      
      const handleSubmit = (e) => {
        e.preventDefault();
        setIsFormSubmitted((prevState) => !prevState);
      };

    const handleSlideChange = (swiper) => {
        setCurrentSlideIndex(swiper.realIndex);
    };

    const macros = dataSearch.Macros
    console.log(Object.entries(dataSearch))
    return (
        <div className='search'>
            <div className='search__top justify--content'>
                <div className='search__top__text'>
                    <h2>Search</h2>
                    <span>
                            Lorem ipsum dolor sit amet, consectetur adipiscing elit. Integer a augue mi.
                            Etiam sed imperdiet ligula, vel elementum velit.
                            Phasellus sodales felis eu condimentum convallis.
                            Suspendisse sodales malesuada iaculis. Mauris molestie placerat ex non malesuada.
                            Curabitur eget sagittis eros. Aliquam aliquam sem non tincidunt volutpat. 
                    </span>
                </div>
                <figure className='searchStats'></figure>
            </div>
            <figure className='line'></figure>
            <div className='research'>
                <form className='justify--content' onSubmit={handleSubmit} >
                    <input
                        type="text"
                        value={gitRef}
                        onChange={handleInputChange}
                        placeholder="Search using commit SHA"
                        className='research__input'
                    />
                    <button type="submit">Search</button>
                </form>
            </div>
            <div className='search__macro justify--content '>
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
                        {dataSearch.Macros && typeof dataSearch.Macros === 'object' && Object.entries(dataSearch.Macros).map(function (search, index)  {
                            return (
                                <SwiperSlide key={index} >
                                    <SearchMacro data={search}/>
                                </SwiperSlide>
                            )
                        })}
                        
                    </Swiper>
                </div>             
                
            </div>
        </div>
    );
};

export default Search;