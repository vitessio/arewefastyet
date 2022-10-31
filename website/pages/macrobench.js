import styles from '../styles/Home.module.css'
import stylesStatus from '../styles/Status.module.css'
import Header from "../components/header";
import Footer from "../components/footer";
import Table from 'react-bootstrap/Table';
import Badge from 'react-bootstrap/Badge';
import moment from "moment";
import { useState, useEffect } from 'react'
import Waiter from "./waiter";

import Button from 'react-bootstrap/Button';
import ButtonGroup from 'react-bootstrap/ButtonGroup';
import Dropdown from 'react-bootstrap/Dropdown';
import DropdownButton from 'react-bootstrap/DropdownButton';
import ButtonToolbar from 'react-bootstrap/ButtonToolbar';

export default function Macrobench() {
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
    }, [])

    if (isVitessRefsLoading || !vitessRefs) {
        return <Waiter />
    }

    return (
        <div>
            <Header />
            <div className={styles.container}>
                <div className={stylesStatus.card}>
                    <h4 className={stylesStatus.h4}>Compare Macrobenchmarks</h4>
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
                        <Button>Compare</Button>
                    </ButtonGroup>
                </div>
                <div className={stylesStatus.card}>
                    <p>Comparing {fromRef.name} with {toRef.name}</p>
                </div>
            </div>
            <Footer />
        </div>
    )
}
