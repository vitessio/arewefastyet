import useGlobalContext from "../contexts/GlobalContext";

export default function useModal() {
  const { modal, setModal } = useGlobalContext();

  function hide() {
    setModal(null);
  }

  function show(element) {
    setModal(element);
  }

  return { element: modal, show, hide };
}
