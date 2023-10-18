import { useEffect } from "react";

export default function useClickOutside(ref, callback) {
  function handleClick(e) {
    if (ref.current && !ref.current.contains(e.target)) {
      callback();
    }
  }

  useEffect(() => {
    document.addEventListener("click", handleClick);

    return () => {
      document.removeEventListener("click", handleClick);
    };
  });
}
