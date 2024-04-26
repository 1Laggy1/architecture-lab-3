export default function moveButtonWrapper() {
  function moveButton() {
    const firstArg = document.getElementById("firstArg").value;
    const secondArg = document.getElementById("secondArg").value;
    if (firstArg <= -1 || firstArg >= 1 || secondArg <= -1 || secondArg >= 1) {
      console.log("wrong argumets");
      return;
    }
    let text;
    text = `move ${firstArg} ${secondArg}`;
    console.log(text);
    fetch("http://localhost:17000", {
      method: "POST",
      body: `${text}`,
      headers: {
        "Content-Type": "text/plain",
      },
    });
  }
  document
    .getElementById("submitButtonMove")
    .addEventListener("click", moveButton);
}
