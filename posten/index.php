<?php
require('../session/session.php');
if(isset($login_session) && $_SESSION['rolle'] == 1) :
?>
<!DOCTYPE html>
<html>
<head>
<title>Georgslauf</title>
<meta name="viewport" content="width=device-width, initial-scale=1"/>
<meta http-equiv="Content-Type" content="text/html; charset=iso-8859-1"/>
<link rel="shortcut icon" type="image/x-icon" href="../res/fav.ico">
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
	<div data-role="page" data-theme="a" id="georgslauf">
		<div data-role="header" data-position="fixed" data-tap-toggle="false" align="center">
			<table border="0" width="100%">
			<tr>
				<td style="width:5%;" align="center"><img src="../res/logo.png" onclick="window.location.href='/'" style="max-height:70px;"/></td>
				<td style="width:90%;">
					<div data-role="navbar" id="menu">
						<ul>
							<li><a target="posten.php" class="link ui-btn-active" style="height:30px; text-shadow:none;">Bewertungen</a></li>
							<li><a target="info.php" class="link" style="height:30px; text-shadow:none;">Posteninfo</a></li>
						</ul>
					</div>
				</td>
			</tr>
			</table>
		</div>
		<div data-role="content" id="home" style="max-width:1200px; margin-left:auto; margin-right:auto;">
		</div>
	</div>
</body>
<noscript>
    <style type="text/css">
        .pagecontainer {display:none;}
    </style>
    <div class="noscriptmsg">
    <strong>Bitte aktiviere JavaScript. Sogar der Georgslauf kommt nicht ohne aus :(</strong>
    </div>
</noscript>
</html>
<?php endif; ?>
