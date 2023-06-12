import React from 'react';

import '../Macro/macro.css'
const Macro = () => {
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
        </div>
    );
};

export default Macro;