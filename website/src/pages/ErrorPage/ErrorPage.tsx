import { Link } from "react-router-dom";

export default function ErrorPage() {
  return (
    <section className="h-screen flex flex-col justify-center items-center gap-y-8">
      <h1 className="text-primary font-bold text-5xl">404</h1>
      <div className="w-[20vw] aspect-square translate-x-[11%]">
        <img src="/404.png" alt="error" />
      </div>
      <h2>We could not find the page you are looking for</h2>
      <Link
        className="bg-primary rounded-lg px-10 py-2 duration-300 hover:scale-105 hover:saturate-200 hover:brightness-150"
        to="/"
      >
        Go Back Home
      </Link>
    </section>
  );
}
