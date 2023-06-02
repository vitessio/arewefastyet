import React from 'react';
import { Routes,Route } from 'react-router-dom';


import Home from './Home/Home';
import Status from './Status/Status';
import Error from '../utils_/Error';
import Layout from '../pages/Layout'



const PublicRoute = () => {
    return (
        <Routes>
        <Route element={<Layout/>}>
            <Route index element={<Home/>}/>

            <Route path='/home' element={<Home/>}/>
            <Route path='/status' element={<Status/>}/>
            
            <Route path='*' element={<Error/>}/>

         </Route> 
     </Routes>
    );
};

export default PublicRoute;