import { useState, useEffect } from "react";

export default function useLocalStorageState(key, defaultValue) {
  const [angka, setAngka] = useState(() => {
    const saved = localStorage.getItem(key);
    return saved ? parseInt(saved) : defaultValue;
  });

  useEffect(() => {
    localStorage.setItem(key, JSON.stringify(angka));
  }, [key, angka]);

  return { angka, setAngka };
}
