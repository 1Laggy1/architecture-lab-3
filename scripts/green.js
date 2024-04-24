export default function wrapperSetGreenBgColor() {
  function setGreenBgColor() {
    fetch("http://localhost:17000", {
      method: "POST",
      body: "green",
      headers: {
        "Content-Type": "text/plain",
      },
    });
  }
  document.querySelector(".green").addEventListener("click", setGreenBgColor);
}
