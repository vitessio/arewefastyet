import { useState, useEffect } from "react";

export default function usePromisedState<T, I = undefined>(value: Promise<T>, initial?: I) {
  const [state, setState] = useState<T | I | undefined>(initial);

  useEffect(() => {
    let isMounted = true;

    value.then((resolvedValue) => {
      if (isMounted) {
        setState(resolvedValue);
      }
    });

    return () => {
      isMounted = false;
    };
  }, []);

  return [state, setState] as T extends Promise<any>
    ? [T | undefined, (value: T | undefined) => void]
    : [T, (value: T) => void];
}
