import styles from '../styles/Home.module.css'
import stylesFooter from '../styles/Footer.module.css'

export default function Footer() {
    return (
        <div className={styles.container}>
            <footer className={stylesFooter.footer}>
                <p>Copyright &copy; 2022</p>
            </footer>
        </div>
    )
}
