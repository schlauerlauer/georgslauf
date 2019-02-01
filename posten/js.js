$(document).ready(function() {
	load('posten.php');

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

	$(document).on('change', '.punkte', function() {
		var id = $(this).attr('id');
		var p = $(this).val();
		if (p > 100) p = 100;
		else if (p < 0) p = 0;
		$(this).val(p);
		$.post('update.php', { p: p, id: id }, function(data) {
			if(data == "ok") $("#"+id).css("color", "green");
			else $("#"+id).css("color", "red");
		});
	});
});
