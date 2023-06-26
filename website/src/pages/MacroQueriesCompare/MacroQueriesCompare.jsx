import React from 'react';
import { useLocation } from 'react-router-dom';

import '../MacroQueriesCompare/macroQueriesCompare.css'
const MacroQueriesCompare = () => {
    let { state } = useLocation()
    console.log(state)
    return (
        <div>
            
        </div>
    );
};

export default MacroQueriesCompare;