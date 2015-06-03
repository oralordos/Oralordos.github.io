function changeImage(e) {
	'use strict';
	var image = document.querySelector('#magic8');
	image.style.backgroundImage = 'url("images/8ball/' + Math.ceil(Math.random() * 20) + '.png")';
}

document.querySelector('button').addEventListener('click', changeImage);