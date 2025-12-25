function myFunc() {
    var po = ["apple", "banana", "cherry"];
    var text = "";
    text = "<ul>"
    po.forEach(function(item) {
        text += "<li>" + item + "</li>";
    });
    text += "</ul>";

    document.getElementById("content").innerHTML = text;
}