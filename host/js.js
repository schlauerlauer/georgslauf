$(document).ready(function() {
	load('home.php');

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

	$(document).on('click', '.color', function() {
		var id = $(this).attr('id');
		var select = "1";
		$.post('update.php', { select: select, id: id }, function(data) {
			alertify.alert('Postenfarbe', "Neu laden nach der Auswahl nötig.<br><br>" + data);
		});
	});

	$(document).on('click', '.select', function() {
		var id = $(this).attr('posten');
		var color = $(this).attr('id');
		$(this).attr('id');
		$.post('update.php', { color: color, id: id }, function(data) {
			alertify.alert().close();
		});
	});

	$(document).on('click', '.punkte', function() {
		var punkte = $(this).attr('id');
		$.post('delete.php', { punkte: punkte }, function(data) {
			if(data == "ok") alertify.success("Eintrag entfernt");
			else alertify.error("Etwas ist schiefgelaufen");
		});
	});

	$(document).on('click', '.gkurz', function() {
		var id = $(this).attr('id');
		var kurz = $(this).attr('kurz');
		alertify.prompt( 'Gruppenkürzel ändern', 'Gruppe ' + id, kurz
               , function(evt, value) {
								 $.post('update.php', { kurz: value, id: id }, function(data) {
						       if(data == "ok") alertify.success("Gespeichert");
						       else alertify.error("Etwas ist schiefgelaufen");
						     });
							  }
               , function() {});
	});

	$(document).on('change', '.kurz', function() {
    var id = $(this).attr('id');
    var k = $(this).val();
    $.post('update.php', { k: k, id: id }, function(data) {
      if(data == "ok") $("#"+id).css("color", "green");
      else $("#"+id).css("color", "red");
    });
  });

	$(document).on('click', '.host', function() {
		var site = "files/" + $(this).attr('host') + ".php";
		$.post(site, { }, function(data) {
			$('#content').html(data);
			$('#content').enhanceWithin();
			window.location="#content";
		});
	});

	$(document).on('click', '#copy', function() {
		var copyText = document.querySelector("#input");
		copyText.select();
		document.execCommand("copy");
	});

	$(document).on('click', '.confirm', function() {
		var site = "files/" + $(this).attr('host') + ".php";
		alertify.confirm("<span style=\"color:red;\">ACHTUNG! Bitte bestätigen</span>", "Wirklich " + $(this).text() + "?", function() {
			$.post(site, { }, function(data) {
				alertify.error(data)
			});
		}, function() {
			alertify.error("AAAaaahA!aah!!");
			alertify.error("AAAaaahAAAAaah!!");
			alertify.warning("AAAAA1!1!Aah!!");
			alertify.error("AAAaahAAAah!!");
			alertify.error("AAAaaa!h!AAAAAaah!!");
			alertify.warning("AAAaahAAAAaah!!");
			alertify.error("AAAaaahAah!!");
			alertify.warning("AAAaaah1Aaah!!");
			alertify.error("AAAahAAAAah!!");
			alertify.warning("AAAaa1ahAaah!!");
			alertify.success(".");
			alertify.success("..");
			alertify.success("...");
			alertify.success("Nichts passiert.");
			alertify.success(";)");
		});
	});

	$(document).on('change', '.coor', function() {
		var xy = $(this).attr('id');
		var posten = $(this).attr('kurz')
		var coor = $(this).val();
		$.post('update.php', { xy: xy, posten: posten, coor: coor }, function(data) {
			if(data == "ok") {}
			else alertify.error(data);
		});
	});

	$(document).on('change', '.startG', function() {
		var kurzel = $(this).attr('kurz');
		var startG = $(this).val();
		$.post('update.php', { startG: startG, kurzel: kurzel }, function(data) {
			if(data == "ok") {}
			else alertify.error(data);
		});
	});

	$(document).on('change', '#auswahl', function() {
		if($(this).attr('type') == 'gruppe') var target = 'files/gruppenPunkte.php';
		else var target = 'files/postenPunkte.php';
		var gruppe = $("#auswahl option:selected").val();
		$.post(target, { gruppe: gruppe }, function(data) {
			$('#punkte').html(data);
			$('#punkte').enhanceWithin();
		})
	})

	$(document).on('change', '.ppunkte', function() {
		var punkte = $(this).val();
		var an = $(this).attr('an');
		var von = $(this).attr('von');
		$.post('punkte.php', { punkte: punkte, von: von, an: an }, function(data) {
			if(data == "ok") {}
			else alertify.error("Etwas ist schiefgelaufen");
		});
	});

});
