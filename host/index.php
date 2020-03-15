<?php
require('../session/session.php');
if(isset($login_session) && $_SESSION['rolle'] >= 3) :
?>
<!DOCTYPE html>
<html>
<head>
<title>Georgslauf</title>
<meta name="viewport" content="width=device-width, initial-scale=1"/>
<meta http-equiv="Content-Type" content="text/html; charset=iso-8859-1"/>
<link rel="stylesheet" href="/js/css/jqm.min.css"/>
<script src="/js/jq.min.js"></script>
<script src="/js/jqm.min.js"></script>
<script src="/js/alertify.min.js"></script>
<link rel="stylesheet" href="/js/css/alertify.min.css" />
<link rel="stylesheet" href="/js/css/themes/default.min.css" />
<meta name="mobile-web-app-capable" content="yes"/>
<script src="js.js"></script>
</head>
<body>
	<div data-role="page" data-theme="b" id="georgslauf">
		<div data-role="content" id="home" style="max-width:1200px; margin-left:auto; margin-right:auto;">
		</div>
	</div>
	<div data-role="footer" data-position="fixed" data-tap-toggle="false" style="overflow:hidden; margin: 0px;">
		<h2><img src="/res/logo.png" style="max-height:30px;"/></h2>
	</div>
</body>
</html>
<?php else : ?>
	Keine Rechte
<?php endif; ?>
