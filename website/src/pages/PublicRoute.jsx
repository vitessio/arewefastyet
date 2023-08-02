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
import { Routes,Route } from 'react-router-dom';


import Home from './Home/Home';
import Status from './Status/Status';
import Error from '../utils/Error/Error';
import Layout from '../pages/Layout'
import Macro from './Macro/Macro';
import Micro from './Micro/Micro';
import Search from './Search/Search';
import CRON from './CRON/CRON';
import MacroQueriesCompare from './MacroQueriesCompare/MacroQueriesCompare';
import Compare from './Compare/Compare';
import PR from './PR/PR';
import SinglePR from './SinglePR/SinglePR';

const PublicRoute = () => {
    return (
        <Routes>
        <Route element={<Layout/>}>
            <Route index element={<Home/>}/>

            <Route path='/home' element={<Home/>}/>
            <Route path='/status' element={<Status/>}/>
            <Route path='/cron' element={<CRON/>}/>
            <Route path='/search' element={<Search/>}/>
            <Route path='/compare' element={<Compare/>}/>
            <Route path='/macro' element={<Macro/>}/>
            <Route path='/macrobench/queries/compare' element={<MacroQueriesCompare/>}/>
            <Route path='/micro' element={<Micro/>}/>
            <Route path='/pr' element={<PR/>}/>
            <Route path='/pr/:pull_nb' element={<SinglePR/>}/>


            <Route path='*' element={<Error/>}/>

         </Route> 
     </Routes>
    );
};

export default PublicRoute;