var timeLeft = null;
var timerId = null;

function countdown() {
    timeLeft--;
    setTimer();
    randomizeHitButton();
}

function randomizeHitButton() {
    var button = document.querySelector('#hit');
    button.style.left = Math.floor(Math.random() * (window.innerWidth - button.clientWidth)) + 'px';
    button.style.top = Math.floor(Math.random() * (window.innerHeight - button.clientHeight)) + 'px';
}

function setTimer() {
    var extraS;
    if (timeLeft === 1) {
        extraS = '';
    }
    else {
        extraS = 's';
    }
    document.querySelector('#timer').innerHTML = timeLeft + ' second' + extraS;

    if (timeLeft <= 5) {
        document.querySelector('#message').innerHTML = 'You are all going to die!';
    }
    else if (timeLeft <= 10) {
        document.querySelector('#message').innerHTML = 'Bad things will happen when time runs out!';
    }
    else if (timeLeft <= 15) {
        document.querySelector('#message').innerHTML = 'Time is running low!';
    }
    else if (timeLeft <= 20) {
        document.querySelector('#message').innerHTML = 'This is going to be dangerous!';
    }
    else if (timeLeft <= 25) {
        document.querySelector('#message').innerHTML = 'The button needs to be pressed!';
    }


    if (timeLeft === 0) {
        document.querySelector('#failure').style.visibility = 'visible';
        abort();
    }
}

function hitButton() {
    timeLeft = 108;
    if (timerId === null) {
        timerId = setInterval(countdown, 1000);
    }
    setTimer();
}

function abort() {
    if (timerId !== null) {
        clearInterval(timerId);
        timerId = null;
    }
    timeLeft = null;
    document.querySelector('#timer').innerHTML = '';
    document.querySelector('#message').innerHTML = '';
}

document.querySelector('#hit').addEventListener('click', hitButton);
document.querySelector('#abort').addEventListener('click', abort);
