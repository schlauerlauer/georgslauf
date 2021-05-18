$(document).ready(function() {
	load('/stamm/stamm.php');

	function load(seite) {
		$('#home').load(seite, function() {
			window.scrollTo(0, 0);
			$('#home').enhanceWithin();
		});
	}

	/* $(document).on('input', '.scoretext', function() {
		var box = $(this);
		var points = box.val();
		box.css('color','red');
		var posten = box.attr('id');
		$.post('set.php', { posten: posten, points: points }, function(data) {
			if (data == points) {
				box.css('color','green');
			}
			else alertify.error("Etwas ist schief gelaufen");
		});
	}); */

	$(document).on('click', '.link', function() {
		var content = $(this).attr('target');
		load(content);
	});

	$(document).on('click', '.map', function() {
		alertify.alert('Posten '+$(this).attr('p')+' Karte', '<iframe frameborder="0" height="460" width="548" src="https://www.google.com/maps/embed/v1/place?q='+$(this).attr('x')+'%2C%20'+$(this).attr('y')+'&amp;key=AIzaSyAp4KBucaAYgMi9WhcelC6g74MAu5iFh_w"></iframe>').set('padding',false);
	});

	$(document).on('click', '.help', function() {
		var help = $(this).attr('help');
		$.post('help.php', { help: help }, function(data) {
			alertify.alert("Hilfe zu den " + help, data);
		});
	});

	$(document).on('click', '.save', function() {
		var gorp = $(this).attr('id');
		if (gorp == "g_save") {
			var gname = $("#g_name").val();
			var gval = $("#g_anzahl").val();
			var gveg = $("#g_veggie").val();
			var gnum = $("#g_kontakt").val();
			var stid = $("#g_form input[type='radio']:checked").val();
			if (gname != "" && gval != "" && stid != null) {
				$.post('insert.php', { gorp: gorp, gname: gname, gval: gval, stid: stid, gveg: gveg, gnum: gnum }, function(data) {
					if(data == 1) alertify.success("Gruppe " + gname + " angemeldet");
					else alertify.error("Etwas ist schiefgelaufen");
				});
			} else alertify.error("Bitte benötigte Angaben eintragen *");
		} else if (gorp == "p_save") {
			var pname = $("#p_name").val();
			var pdesc = $("#p_desc").val();
			var pkont = $("#p_kont").val();
			var pval = $("#p_anzahl").val();
			var pveg = 0;
			var pmat = $("#p_mat").val();
			var port = $("#p_ort").val();
			var	psonst = $("#p_sonst").val();
			var pid = 2;
			var kid = $("#k_form input[type='radio']:checked").val();
			if (pname != "" && pdesc != "" && pkont != "" && pval != "" && pveg != "" && pid != null && kid != null) {
				$.post('insert.php', { gorp: gorp, pname: pname, pdesc: pdesc, pkont: pkont, pval: pval, pveg: pveg, pmat: pmat, port: port, psonst: psonst, kid: kid, pid : pid }, function(data) {
					if (data == 1) alertify.success("Posten " + pname + " angemeldet");
					else if(data==2) alertify.error("Diese Kategorie für " + pid + " ist voll");
					else alertify.error("Etwas ist schiefgelaufen");
				});
			} else alertify.error("Bitte benötigte Angaben eintragen *");
		}
	});

	$(document).on('click', '.delete', function() {
		var pid = $(this).attr('id');
		var type = $(this).attr('type');
		alertify.confirm('Abmelden', 'Soll ' + $(this).attr('name') + ' entfernt werden?', function() {
				$.post('delete.php', { pid: pid, type: type }, function(data) {
					if(data == 'ok') alertify.success("Eintrag entfernt");
					if(data == 'nein') alertify.error("Abmeldung nur noch per Email möglich");
					else alertify.error("Etwas ist schiefgelaufen");
					if(type == "p") load('/stamm/posten.php');
					else load('/stamm/gruppen.php');
				});
			}, function() {});
	});
});
