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

function Counter() {
  new pn = function() { // 箭头函数：这里的this指向定义时所在对象（即Counter实例）
    this.count = 0;
    setInterval(() => {
      this.count++; // 正确！这里的this指向Counter实例
    }, 1000);
  }
}
