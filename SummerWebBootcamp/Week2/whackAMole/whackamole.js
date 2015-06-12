var points = 0;
var molePopup = null;
var size = 2;
var interval = 10;

function createMoles(x, y) {
    var game = document.querySelector('#game');
    game.style.width = x * 175 + 8 * x + 'px';
    game.style.height = y * 150 + 8 * y + 'px';
    for (var i = 0; i < y; i++) {
        var newRow = document.createElement('div');
        newRow.className = 'game-row';
        for (var j = 0; j < x; j++) {
            var newMole = document.createElement('img');
            newMole.className = 'mole';
            newMole.src = 'whackamole.png';
            newMole.dataset.timerId = 'null';
            newMole.addEventListener('click', onClick);
            newRow.appendChild(newMole);
        }
        game.appendChild(newRow);
    }
    molePopup = setTimeout(randomMole, 2000);
}

function startBoss() {
    clearTimeout(molePopup);
    molePopup = null;
    var moles = document.querySelectorAll('.mole');
    for (var i = 0; i < moles.length; i++) {
        moles[i].removeEventListener('click', onClick);
        moles[i].parentNode.removeChild(moles[i]);
    }
    var rows = document.querySelectorAll('.game-row');
    for (i = 0; i < rows.length; i++) {
        rows[i].parentNode.removeChild(rows[i]);
    }
    var boss = document.createElement('img');
    boss.className = 'boss';
    boss.src = 'whackamole.png';
    var game = document.querySelector('#game');
    boss.dataset.intervalTimer = setInterval(function () {
        boss.style.left = Math.floor(Math.random() * (game.clientWidth - boss.clientWidth)) + 'px';
        boss.style.top = Math.floor(Math.random() * (game.clientHeight - boss.clientHeight)) + 'px';
    }, 1000);
    boss.addEventListener('mouseenter', function () {
        if (boss.x < (game.clientWidth - boss.clientWidth) / 2) {
            boss.style.left = (game.clientWidth - boss.clientWidth) + 'px';
        }
        else {
            boss.style.left = '0';
        }
        if (boss.y < (game.clientHeight - boss.clientHeight) / 2) {
            boss.style.top = (game.clientHeight - boss.clientHeight) + 'px';
        }
        else {
            boss.style.top = '0';
        }
    });
    var health = 3;
    boss.addEventListener('click', function () {
        health--;
        boss.classList.add('hit');
        setTimeout(function () {
            boss.classList.remove('hit');
        }, 250);
        playHit();
        if (health <= 0) {
            clearInterval(boss.dataset.intervalTimer);
            boss.parentNode.removeChild(boss);
            points += 40;
            updateScores();
            createMoles(size, size);
            size++;
        }
    });
    game.appendChild(boss);
}

function updateScores() {
    document.querySelector('#points').innerHTML = 'Points: ' + points;
}

function playHit() {
    var newSound = document.createElement('audio');
    newSound.src = 'Socapex%20-%20big%20punch.mp3';
    newSound.play();
    document.body.appendChild(newSound);
    setTimeout(function () {
        document.body.removeChild(newSound);
    }, 1000);
}

function onClick(e) {
    if (e.target.dataset.timerId !== 'null') {
        points++;
        if (points % interval === 0) {
            startBoss();
        }
        updateScores();
        clearTimeout(e.target.dataset.timerId);
        e.target.dataset.timerId = null;
        e.target.classList.add('hit');
        playHit();
        setTimeout(function () {
            e.target.style.opacity = '';
            e.target.classList.remove('hit');
        }, 250);
    }
}

function randomMole() {
    var moles = document.querySelectorAll('.mole');
    var whichMole;
    var isDownMole = false;
    for (var i = 0; i < moles.length; i++) {
        if (moles[i].style.opacity === '') {
            isDownMole = true;
            break;
        }
    }
    if (isDownMole) {
        do {
            whichMole = moles[Math.floor(Math.random() * moles.length)];
        } while (whichMole.style.opacity !== '');
        whichMole.style.opacity = '1';
        var delay = Math.floor(Math.random() * 4000 + 1000);
        whichMole.dataset.timerId = setTimeout(function () {
            whichMole.style.opacity = '';
            updateScores();
            whichMole.dataset.timerId = 'null';
        }, delay);
    }
    molePopup = setTimeout(randomMole, Math.floor(Math.random() * 2000 + 500));
}

function onLoad() {
    createMoles(size, size);
    size++;
}

window.addEventListener('load', onLoad);