$(document).ready(function() {
	load('test_x.php');

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

  $(document).on('change', '.punkte', function() {
    var id = $(this).attr('id');
    var p = $(this).val();
    if (p > 100) p = 100;
    else if (p < 0) p = 0;
    $(this).val(p);
    $.post('test_points.php', { p: p, id: id }, function(data) {
      if(data == "ok") $("#"+id).css("color", "green");
      else $("#"+id).css("color", "red");
    });
  });

	var simplemde = new SimpleMDE({ element: document.getElementById("smde") });
	simplemde.value("Testtext");

});
