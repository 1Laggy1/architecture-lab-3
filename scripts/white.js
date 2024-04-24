export default function wrapperSetWhiteBgColor() {
  function setWhiteBgColor() {
    fetch("http://localhost:17000", {
      method: "POST",
      body: "white",
      headers: {
        "Content-Type": "text/plain",
      },
    });
  }
  document.querySelector(".white").addEventListener("click", setWhiteBgColor);
}
