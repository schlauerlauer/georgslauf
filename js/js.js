$(document).ready(function() {
	load('../info.php');
	countdown();

	function load(seite) {
		$('#home').load(seite, function() {
			window.scrollTo(0, 0);
			$('#home').enhanceWithin();
		});
	}
	$(document).on('click', '.link', function() {
		var content = $(this).attr('target');
		load(content);
	});

	$(document).on('click', '.map', function() {
		alertify.alert('Posten '+$(this).attr('p')+' Karte', '<iframe frameborder="0" height="460" width="548" src="https://www.google.com/maps/embed/v1/place?q='+$(this).attr('x')+'%2C%20'+$(this).attr('y')+'&amp;key=AIzaSyAp4KBucaAYgMi9WhcelC6g74MAu5iFh_w"></iframe>').set('padding',false);
	});
});

function login() {
	var un = $('#username').val();
	var pw = $('#password').val();
	$.post('/session/login.php', { un: un, pw: pw }, function(data) {
		if(data == "2") window.location.href = '/stamm';
		else if(data == "1") window.location.href = '/posten';
		else if(data == "3") window.location.href = '/host';
		else alertify.error("Benutzername oder Passwort falsch");
	});
	}

function countdown() {
	var countDownDate = new Date("Apr 21, 2018 08:00:00").getTime();
	var x = setInterval(function() {
  	var now = new Date().getTime();
		var distance = countDownDate - now;

    var days = Math.floor(distance / (1000 * 60 * 60 * 24));
    var hours = Math.floor((distance % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
    var minutes = Math.floor((distance % (1000 * 60 * 60)) / (1000 * 60));
    var seconds = Math.floor((distance % (1000 * 60)) / 1000);
		var min = " Minute";
		if(minutes != 1 ) min += "n";
    document.getElementById("clock").innerHTML = days + " Tage " + hours + " Stunden "
    + minutes + min;

    if (distance < 0) {
        clearInterval(x);
        document.getElementById("clock").innerHTML = "";
    }
	}, 1000);
}
