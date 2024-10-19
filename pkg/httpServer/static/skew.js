



// body.style.transform = "skewY(-1.75deg)";

window.onload = () => {
    const body = document.getElementsByTagName("body")[0]
    x = ((Math.random() * 10) - 5).toFixed(2)
    y = ((Math.random() * 10) - 5).toFixed(2)
    console.log(x,y)
    
    body.style.transform = "skewY(-1.75deg)";
}

