export default function wrapperUpdate() {
  function update() {
    fetch("http://localhost:17000", {
      method: "POST",
      body: "update",
      headers: {
        "Content-Type": "text/plain",
      },
    });
  }
  document.querySelector(".update").addEventListener("click", update);
}
