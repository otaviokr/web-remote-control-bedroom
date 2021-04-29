// This function will update the LEDs on the image to reflect the selected color.
function changeColor(AreaID) {
    console.log("This is changeColor with AreaID: " + AreaID);
    var obj = document.getElementById("parent_svg").contentDocument.getElementById(AreaID);
    console.log(obj);
    var newColor = document.getElementById("colorpicker").value;
    console.log(newColor);
    var newBright = document.getElementById("bright").value;
    obj.style.fill = newColor;
    obj.style.stroke = newColor;
    document.getElementById("input_" + AreaID).value = newColor;
    console.log(document.getElementById("input_" + AreaID));
    console.log(document.getElementById("input_" + AreaID).value);
    document.getElementById("input_" + AreaID + "b").value = newBright;
    console.log(document.getElementById("input_" + AreaID + "b"));
    console.log(document.getElementById("input_" + AreaID + "b").value);
}

// You have to convert the form, otherwise, the data is sent as multipart/form-data
// and it will not be parsable by golang.
function urlencodeFormData(fd) {
    var params = new URLSearchParams();
    for(var pair of fd.entries()) {
        typeof pair[1] == 'string' && params.append(pair[0], pair[1]);
    }
    return params.toString();
}

// Async call to light the LEDs on Blinkt.
function submitColors() {
    var f = document.getElementById("form");
    var formData = urlencodeFormData(new FormData(f));
    const xhr = new XMLHttpRequest();
    xhr.open("POST", "/.netlify/functions/blinkt");
    xhr.setRequestHeader('Content-type', 'application/x-www-form-urlencoded');
    xhr.send(formData);
}