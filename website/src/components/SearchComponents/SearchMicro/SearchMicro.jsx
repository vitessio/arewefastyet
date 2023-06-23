import React, {useState} from 'react';

import './searchMicro.css'

import { closeTables, openTables } from '../../../utils/Utils';

const SearchMicro = ({data, className}) => {
    const [maxHeight, setMaxHeight] = useState(closeTables);

    const handleClick = () => {
        if (maxHeight === closeTables) {
            setMaxHeight(openTables);
          } else {
            setMaxHeight(closeTables);
        }
    };

    return (
        <div className={`microbench  ${className}`} style={{ maxHeight: `${maxHeight}px` }}>
            <div className='space--between justify--content align--center'>
                <span className='width--12em'>{data[1].PkgName}</span>
                <span className='width--14em name'>{data[1].Name}</span>
                <span className='width--18em hiddenMobile'>{data[1].Result.Ops.toFixed(0)}</span>
                <span className='width--18em hiddenTablet'>{data[1].Result.NSPerOp.toFixed(0)}</span>
                <div className='width--6em'><i className="fa-solid fa-circle-info" onClick={handleClick}></i></div>
            </div>
            <div className='search__microbench__bottom flex--column'>
                <div className='hiddenDesktop microbenchMore'>
                    <span className='width--18em'>Number of Iterations</span>
                    <span className='width--12em'>{data[1].Result.Ops.toFixed(0)}</span>
                </div>
                <div className=' hiddenDesktop microbenchMore'>
                    <span className='width--18em'>Time/op</span>
                    <span className='width--12em'>{data[1].Result.NSPerOp.toFixed(0)}</span>
                </div>
                <div className='flex microbenchMore'>
                    <span className='width--18em'>Bytes/op</span>
                    <span className='width--12em'>{data[1].Result.BytesPerOp}</span>
                </div>
                <div className='flex microbenchMore'>
                    <span className='width--18em'>Megabytess/s</span>
                    <span className='width--12em'>{data[1].Result.MBPerSec}</span>
                </div>
                <div className='flex microbenchMore'>
                    <span className='width--18em'>Allocations/op</span>
                    <span className='width--12em'>{data[1].Result.AllocsPerOp}</span>
                </div>
            </div>
        </div>
    );
};

export default SearchMicro;