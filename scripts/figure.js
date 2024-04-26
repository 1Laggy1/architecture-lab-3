export default function figureButtonWrapper() {
  function figureButton() {
    const firstArg = document.getElementById("firstArgFigure").value;
    const secondArg = document.getElementById("secondArgFigure").value;
    console.log("clicked");
    if (!(firstArg <= 0 || firstArg >= 1 || secondArg <= 0 || secondArg >= 0)) {
      console.log("wrong argumets");
      return;
    }
    let text;
    text = `figure ${firstArg} ${secondArg}`;
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
    .getElementById("submitButtonFigure")
    .addEventListener("click", figureButton);
}
