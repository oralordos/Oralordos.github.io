var numMisses = 0;
var numHits = 0;

function createMoles(x, y) {
    var game = document.querySelector('#game');
    for (var i = 0; i < y; i++) {
        var newRow = document.createElement('div');
        newRow.className = 'game-row';
        for (var j = 0; j < x; j++) {
            var newMole = document.createElement('img');
            newMole.className = 'mole';
            newMole.style.opacity = '0.3';
            newMole.src = 'images/whackamole.jpg';
            newMole.dataset.timerId = 'null';
            newRow.appendChild(newMole);
        }
        game.appendChild(newRow);
    }
}

function updateScores() {
    document.querySelector('#numMisses').innerHTML = 'Misses: ' + numMisses;
    document.querySelector('#numHits').innerHTML = 'Hits: ' + numHits;
}

function onClick(e) {
    e.target.style.opacity = '0.3';
    if (e.target.dataset.timerId !== 'null') {
        numHits++;
        updateScores();
        clearTimeout(e.target.dataset.timerId);
        e.target.dataset.timerId = null;
    }
}

function randomMole() {
    var moles = document.querySelectorAll('.mole');
    var whichMole;
    do {
        whichMole = moles[Math.floor(Math.random() * moles.length)];
    } while (whichMole.dataset.timerId !== 'null');
    whichMole.style.opacity = '1';
    var delay = Math.floor(Math.random() * 4000 + 1000);
    whichMole.dataset.timerId = setTimeout(function() {
        whichMole.style.opacity = '0.3';
        numMisses++;
        updateScores();
        whichMole.dataset.timerId = 'null';
    }, delay);
    setTimeout(randomMole, Math.floor(Math.random() * 2000 + 500));
}

function onLoad() {
    createMoles(3, 3);

    var moles = document.querySelectorAll('.mole');
    for (var i = 0; i < moles.length; i++) {
        moles[i].addEventListener('click', onClick);
    }
    setTimeout(randomMole, 2000);
}

window.addEventListener('load', onLoad);