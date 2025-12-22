var text = "Hello World!";
var text2 = "Hello JavaScript!";

var text3 = text2.split(" ");

var html = `<h1>${text}</h1><ul>`;

for (const x of text3) {
    html += `<li>${x}</li>`;
}
html += `</ul>`;


function displayDate() {
    document.getElementById("myText2").innerHTML = html;
}