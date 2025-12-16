function Myfunc(){
    var el = document.getElementById("myText");
    var x = el.textContent !== undefined ? el.textContent : el.value;
    document.getElementById("demo").innerHTML = x;
}