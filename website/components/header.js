import styles from '../styles/Header.module.css'
import Container from 'react-bootstrap/Container';
import Navbar from 'react-bootstrap/Navbar';
import Nav from 'react-bootstrap/Nav';

export default function Header() {
    return (
        <div className={styles.header}>
            <Navbar collapseOnSelect expand="md" bg="dark" variant="dark" fixed="top">
                <Container>
                    <Navbar.Brand href="/">
                        <img
                            alt="Vitess icon"
                            src="https://vitess.io/img/logos/vitess.png"
                            width="30"
                            height="30"
                            className="d-inline-block align-top"
                        />{' '}
                        arewefastyet
                    </Navbar.Brand>
                    <Navbar.Toggle aria-controls="responsive-navbar-nav" />
                    <Navbar.Collapse id="responsive-navbar-nav">
                        <Nav className="me-auto">
                            <Nav.Link href="/#home">Home</Nav.Link>
                            <Nav.Link href="/status">Status</Nav.Link>
                            <Nav.Link href="/#compare">Compare</Nav.Link>
                            <Nav.Link href="/#compare">Search</Nav.Link>
                            <Nav.Link href="/#micro">Micro</Nav.Link>
                            <Nav.Link href="/macrobench">Macro</Nav.Link>
                        </Nav>
                    </Navbar.Collapse>
                </Container>
            </Navbar>
        </div>
    )
}
