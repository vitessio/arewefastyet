import styles from '../styles/Home.module.css'
import stylesCard from '../styles/Card.module.css'
import Header from "../components/header";
import Footer from "../components/footer";
import Waiter from "./waiter";
import MicrobenchCmp from "../components/microbench_cmp";
import { useState, useEffect } from 'react'

import ButtonGroup from 'react-bootstrap/ButtonGroup';
import Dropdown from 'react-bootstrap/Dropdown';
import DropdownButton from 'react-bootstrap/DropdownButton';


export default function Microbench(props) {
    const [fromRef, setFromRef] = useState(null);
    const [toRef, setToRef] = useState(null);

    const [vitessRefs, setVitessRefs] = useState(null);
    const [isVitessRefsLoading, setVitessRefsLoading] = useState(true)

    useEffect(() => {
        fetch('http://localhost:9090/api/vitess/refs')
            .then((res) => res.json())
            .then((data) => {
                setVitessRefs(data)
                setVitessRefsLoading(false)
                setFromRef(data[1])
                setToRef(data[0])
            })
    }, [props])

    if (isVitessRefsLoading || !vitessRefs) {
        return <Waiter />
    }

    return (
        <div>
            <Header />
            <div className={styles.container}>
                <div className={stylesCard.card}>
                    <h4 className={stylesCard.h4}>Compare Microbenchmarks</h4>
                    <ButtonGroup>
                        <DropdownButton variant="light" as={ButtonGroup} title={fromRef.name} id="bg-nested-dropdown-from">
                            {vitessRefs.map((item,index)=>{
                                return <Dropdown.Item key={item.name} onClick={(e) => setFromRef(item)}>{item.name}</Dropdown.Item>
                            })}
                        </DropdownButton>
                        <DropdownButton variant="light" as={ButtonGroup} title={toRef.name} id="bg-nested-dropdown-to">
                            {vitessRefs.map((item,index)=>{
                                return <Dropdown.Item key={item.name} onClick={(e) => setToRef(item)}>{item.name}</Dropdown.Item>
                            })}
                        </DropdownButton>
                    </ButtonGroup>
                </div>
                <MicrobenchCmp from={fromRef} to={toRef} />
            </div>
            <Footer />
        </div>
    )
}
